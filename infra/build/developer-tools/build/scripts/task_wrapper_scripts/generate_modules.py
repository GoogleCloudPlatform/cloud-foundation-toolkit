#!/usr/bin/env python3

# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import subprocess
import sys
import json

from jinja2 import Environment, FileSystemLoader

AUTOGEN_NOTE = '// This file was automatically generated from a template in {folder}'

class Module(object):
    path = None
    options = {}

    def __init__(self, path, template_options):
        self.path = path
        self.options = template_options

    def template_options(self, base):
        return {k: v for d in [base, self.options] for k, v in d.items()}

DEVNULL_FILE = open(os.devnull, 'w')

def main(argv):
    modules = json.loads(argv[1])
    for module in modules:
        template_folder = module["template_folder"]
        module = Module(module["path"], module["options"])
        env = Environment(
            keep_trailing_newline=True,
            loader=FileSystemLoader(template_folder),
            trim_blocks=True,
            lstrip_blocks=True,
        )
        templates = env.list_templates()

        for template_file in templates:
            template = env.get_template(template_file)
            if template_file.endswith(".tf.tmpl"):
                template_file = template_file.replace(".tf.tmpl", ".tf")
            rendered = template.render(
                module.template_options(
                    {'autogeneration_note': AUTOGEN_NOTE.format(folder=template_folder)}
                )
            )
            with open(os.path.join(module.path, template_file), "w") as f:
                f.write(rendered)
                if template_file.endswith(".tf"):
                    subprocess.call(
                        [
                            "terraform",
                            "fmt",
                            "-write=true",
                            os.path.join(module.path, template_file)
                        ],
                        stdout=DEVNULL_FILE,
                        stderr=subprocess.STDOUT
                    )
                if template_file.endswith(".sh"):
                    os.chmod(os.path.join(module.path, template_file), 0o755)
    DEVNULL_FILE.close()


if __name__ == "__main__":
    main(sys.argv)
