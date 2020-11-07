import logging
import requests
from slugify import slugify
import json
import yaml


class GrafanaDBExporter:
    def __init__(self, url, key):
        self._url = url
        self._headers = {'Authorization': f'Bearer {key}'}

    def _call(self, endpoint):
        response = requests.get(f'{self._url}{endpoint}', headers=self._headers)
        if response.status_code != 200:
            logging.warning(f'Call to {endpoint} failed: {response.status_code} - {response.reason}')
            return None
        return response.json()

    def export_datasources(self, filename):
        datasources = {
            'apiVersion': '1',
            'datasources': []
        }
        for source in self._call('/api/datasources'):
            if (datasource := self._call(f'/api/datasources/{source["id"]}')) is not None:
                datasources['datasources'].append(datasource)
        with open(filename,'w') as f:
            f.write(yaml.dump(datasources))

    def export_dashboards(self, directory):
        for dashboard in self._call('/api/search?query=&/'):
            if (export := self._call(f'/api/dashboards/uid/{dashboard["uid"]}')) is not None:
                with open(f'{directory}/{slugify(dashboard["title"])}.json', 'w') as f:
                    f.write(json.dumps(export['dashboard'], indent=4))

    def get_dashboards_old(self):
        url = f'{self._url}/api/search?query=&/'
        response = requests.get(url, headers=self._headers)
        if response.status_code != 200:
            logging.warning(f'Failed to get dashboards: {response.status_code} - {response.reason}')
            return None
        return [dashboard['uri'] for dashboard in response.json()]

    def export_dashboards_old(self):
        for dashboard in self.get_dashboards_old():
            url = f'{self._url}/api/dashboards/{dashboard}'
            response = requests.get(url, headers=self._headers)
            if response.status_code != 200:
                logging.warning(f'Failed to get dashboard: {response.status_code} - {response.reason}')
                continue
            logging.debug(response.json())
            if response.json()['meta']['provisioned']:
                logging.debug(f'{dashboard} is a provisioned dashboard. Skipping')
                continue
            export = response.json()['dashboard']
            logging.info(f'{dashboard}: {export}')
            with open('dashboards/' + slugify(dashboard) + '.json', 'w') as f:
                f.write(json.dumps(export))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    api_key = 'eyJrIjoiRXVEcTNFQm0zb0JXMWNGNnhEOHJjandFOFptTTlzNUwiLCJuIjoiZ3JhZmFuYV9leHBvcnRlciIsImlkIjoxfQ=='
    exporter = GrafanaDBExporter('http://grafana.192.168.0.11.nip.io', api_key)
    exporter.export_datasources('datasources/datasources.yml')
    exporter.export_dashboards('dashboards')
