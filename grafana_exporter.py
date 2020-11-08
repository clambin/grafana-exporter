import os
from pathlib import Path
import logging
import requests
from slugify import slugify
import json
import yaml
from yaml.representer import SafeRepresenter


class LiteralStr(str):
    pass


def change_style(style, representer):
    def new_representer(dumper, data):
        scalar = representer(dumper, data)
        scalar.style = style
        return scalar
    return new_representer


represent_literal_str = change_style('|', SafeRepresenter.represent_str)
yaml.add_representer(LiteralStr, represent_literal_str)


class GrafanaDBExporter:
    def __init__(self, url, key):
        self._url = url
        self._headers = {'Authorization': f'Bearer {key}'}

    def _call(self, endpoint):
        response = requests.get(f'{self._url}{endpoint}', headers=self._headers)
        if response.status_code != 200:
            logging.warning(f'Call to {endpoint} failed: {response.status_code} - {response.reason}')
            return dict()
        return response.json()

    @staticmethod
    def _build_configmap(files, name, namespace, indent=2):
        configmap = {
            'kind': 'ConfigMap',
            'apiVersion': 'v1',
            'metadata': {
                'name': name,
                'namespace': namespace,
            },
            'data': dict()
        }
        for filename, content in files.items():
            configmap['data'][filename] = LiteralStr(yaml.dump(content, indent=indent))
        return configmap

    @staticmethod
    def _save_yaml(path, content, indent=2):
        Path(os.path.dirname(path)).mkdir(parents=True, exist_ok=True)
        with open(path, 'w') as f:
            f.write(yaml.dump(content, indent=indent))
            f.close()

    @staticmethod
    def _save_json(path, content, indent=4):
        Path(os.path.dirname(path)).mkdir(parents=True, exist_ok=True)
        with open(path, 'w') as f:
            f.write(json.dumps(content, indent=indent))
            f.close()

    def _get_datasources(self):
        datasources = {
            'apiVersion': 1,
            'datasources': []
        }
        for source in self._call('/api/datasources'):
            if (datasource := self._call(f'/api/datasources/{source["id"]}')) is not None:
                datasources['datasources'].append(datasource)
        return datasources

    def export_datasources(self, directory):
        GrafanaDBExporter._save_yaml(os.path.join(directory, 'datasources.yml'), self._get_datasources())

    def export_datasources_configmap(self, directory, name='grafana-provisioning-datasources', namespace='monitoring'):
        GrafanaDBExporter._save_yaml(
            os.path.join(directory, f'{name}.yml'),
            GrafanaDBExporter._build_configmap({'datasources.yml': self._get_datasources()}, name, namespace)
        )

    def _get_dashboard_folders(self, folders):
        ids = dict()
        for folder in self._call('/api/folders'):
            if folder['title'] in folders:
                ids[folder['title']] = folder['id']
        if 'General' in folders:
            ids['General'] = 0
        return ids

    def _get_dashboards(self, folderid):
        dashboards = {}
        query = f'folderIds={folderid}&type=dash-db'
        for dashboard in self._call(f'/api/search?{query}'):
            if (export := self._call(f'/api/dashboards/uid/{dashboard["uid"]}')) is not None:
                if not export['meta']['isFolder']:
                    dashboards[slugify(dashboard["title"])] = export['dashboard']
        return dashboards

    @staticmethod
    def _get_dashboard_provisioning(name):
        return {
            'apiVersion': 1,
            'providers': [{
                'name': name,
                'orgId': 1,
                'folder': '',
                'disableDeletion': False,
                'updataIntervalSeconds': 10,
                'allowUiUpdates': True,
                'options': {
                    'path': '/var/lib/grafana/dashboards',
                    'foldersFromFilesStructure': True
                }
            }]
        }

    def export_dashboards(self, directory, folders):
        GrafanaDBExporter._save_yaml(
            os.path.join(directory, 'dashboards.yml'),
            GrafanaDBExporter._get_dashboard_provisioning('my dashboards')
        )
        for foldername, folderid in self._get_dashboard_folders(folders).items():
            for title, dashboard in self._get_dashboards(folderid).items():
                logging.info(f'Exporting {foldername}/{title}')
                GrafanaDBExporter._save_json(os.path.join(directory, foldername, f'{title}.json'), dashboard)

    def export_dashboards_configmap(self, directory, folders, name='grafana-provisioning-dashboards', namespace='monitoring'):
        configmap = GrafanaDBExporter._build_configmap(
            {'dashboards.yml': GrafanaDBExporter._get_dashboard_provisioning('my dashboards')},
            name, namespace
        )
        GrafanaDBExporter._save_yaml(os.path.join(directory, 'grafana-provisioning-dashboards.yml'), configmap)

        for foldername, folderid in self._get_dashboard_folders(folders).items():
            configmap = GrafanaDBExporter._build_configmap({}, f'grafana-dashboards-{slugify(foldername)}', namespace)
            for title, dashboard in self._get_dashboards(folderid).items():
                logging.info(f'Exporting (configmap) {foldername}/{title}')
                configmap['data'][f'{title}.json'] = LiteralStr(json.dumps(dashboard, indent=2))
            GrafanaDBExporter._save_yaml(
                os.path.join(directory, f'{slugify(configmap["metadata"]["name"])}.yml'),
                configmap
            )


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', datefmt='%Y-%m-%d %H:%M:%S',
                        level=logging.INFO)
    api_key = 'eyJrIjoiRXVEcTNFQm0zb0JXMWNGNnhEOHJjandFOFptTTlzNUwiLCJuIjoiZ3JhZmFuYV9leHBvcnRlciIsImlkIjoxfQ=='
    exporter = GrafanaDBExporter('http://grafana.192.168.0.11.nip.io', api_key)
    exporter.export_datasources('datasources')
    exporter.export_dashboards('dashboards', ['General', 'system', 'media', 'covid'])
    exporter.export_datasources_configmap('k8s')
    exporter.export_dashboards_configmap('k8s', ['General', 'system', 'media', 'covid'])
