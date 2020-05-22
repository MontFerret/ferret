import os

from ruamel.yaml import YAML


def walk_indices(directory):
  for root, dirs, files in os.walk(directory):
    for file in files:

      if file == 'lib.go':
        yield os.path.join(root, file)


def walk_implementations(directory):
  for root, dirs, files in os.walk(directory):
    for file in files:

      if file == 'lib.go':
        continue

      if file.endswith('_test.go'):
        continue

      yield os.path.join(root, file)


def dump_yaml(object, file_name):
  with open(file_name, 'w+') as file:
    YAML().dump(object, file)
