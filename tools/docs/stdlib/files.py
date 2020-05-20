import os


def walk(directory):
  for root, dirs, files in os.walk(directory):
    for file in files:

      if file == 'lib.go':
        continue

      if file.endswith('_test.go'):
        continue

      yield os.path.join(root, file)
