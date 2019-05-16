from os.path import join, dirname, isdir
from os import walk
import yaml

def read_file(path):
  with open(path) as f:
    contents = f.read()
    f.close()
    return yaml.safe_load(contents)

def look_for_schemas_tests(root):
  ret = []

  for (dirpath, dirnames, filenames) in walk(root):
    for file in filenames:
      if (
          file.endswith('.yaml') and
          (file.startswith('invalid_') or file.startswith('valid_'))
      ):
        filename = join(dirpath, file)
        ret.append(
          (file.startswith('valid_'), filename, read_file(filename))
        )

  return ret

def look_for_schemas_dirs_tests():
  template_root = join(dirname(__file__), '..', '..', 'templates')
  ret = []

  for (dirpath, dirnames, filenames) in walk(template_root):
    for dir in dirnames:
      dir_unit = join(dirpath, dir, 'tests', 'schemas')
      if isdir(dir_unit):
        schema = read_file(join(dirpath, dir, dir + '.py.schema'))
        ret.append((dir, schema, look_for_schemas_tests(dir_unit)))

  return ret

