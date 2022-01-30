#!/bin/sh

if [ -z "$GIT_USER" ]; then
	echo Missing GIT_USER env var
	exit 1
fi

if [ -z "$GIT_EMAIL" ]; then
	echo Missing GIT_EMAIL env var
	exit 1
fi

if [ -z "$GIT_FULL_NAME" ]; then
	echo Missing GIT_FULL_NAME env var
	exit 1
fi

if [ -z "$GIT_TOKEN" ]; then
	echo Missing GIT_TOKEN env var
	exit 1
fi

if [ -z "$GRAFANA_API_KEY" ]; then
  echo Missing GRAFANA_API_KEY env var
  exit 1
fi

[ -z "$GRAFANA_URL" ] && GRAFANA_URL=http://grafana.monitoring.svc:3000

git config --global user.email "$GIT_EMAIL" || exit 1
git config --global user.name "$GIT_FULL_NAME" || exit 1

TMPDIR=$(mktemp -d)
git clone https://"$GIT_USER":"$GIT_TOKEN"@github.com/clambin/gitops.git "$TMPDIR" || exit 1

if [ -n "$BRANCH" ]; then
  git checkout "$BRANCH"
fi

cd "$TMPDIR" || exit 1

if [ -n "$OUT_DATASOURCES" ]; then
  /app/grafana-exporter --out "$OUT_DATASOURCES" --url "$GRAFANA_URL" --token "$GRAFANA_API_KEY" datasources || exit 1
fi

if [ -n "$OUT_DASHBOARDS" ]; then
  /app/grafana-exporter --out "$OUT_DASHBOARDS" --url "$GRAFANA_URL" --token "$GRAFANA_API_KEY" dashboard-provisioning || exit 1
  /app/grafana-exporter --out "$OUT_DASHBOARDS" --url "$GRAFANA_URL" --token "$GRAFANA_API_KEY" dashboards --folders="$GRAFANA_FOLDERS" || exit 1
fi

if [ -z "$SKIP_COMMIT" ]; then
  git add -A &&
  git commit -m "Automated grafana export on $(date +'%Y-%m-%d %H:%M:%S')" &&
  git push
  echo "Successfully synced grafana configuration with git"
fi

cd - >/dev/null || exit
rm -rf "$TMPDIR"
