from .conftest import look_for_schemas_dirs_tests
from jsonschema.exceptions import ValidationError
from jsonschema import validate

def test_schemas():
  modules = look_for_schemas_dirs_tests()
  for (module, schema, files) in modules:
    for (isValid, path, data) in files:
      try:
        validate(data, schema)
        if not isValid:
          raise Exception("Validation for {} should have failed".format(path))
      except ValidationError:
        if isValid:
          raise
