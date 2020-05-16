#!/usr/bin/env python3

# Copyright 2020 Google LLC
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
import sys
import argparse
import yaml

import requests
from collections import OrderedDict
from jinja2 import Environment, FileSystemLoader

TERRAFORM_REGISTRY_BASE = "https://registry.terraform.io/v1/"

class IndexItem(yaml.YAMLObject):
  yaml_tag = u'!module'

  def __init__(self, data):
    self.children = {}
    self.data = data

  def name(self):
    return self.data.get("name")

  def url(self):
    if "source" in self.data:
      return self.data.get("source")
    path = self.data.get("path")
    return f"{self.parent.url()}/tree/master/{path}"

  def should_display(self):
    return not self.data.get('exclude', False)

  def description(self):
    return self.data.get("description")

  def add_child_data(self, data):
    child = self.add_child(IndexItem(data))
    child.data = {**child.data, **data}
    return child

  def add_child(self, child):
    if child.name() not in self.children:
      self.children[child.name()] = child
      child.parent = self
    return self.children[child.name()]

  @classmethod
  def to_yaml(cls, representer, node):
    rep = OrderedDict()
    rep_keys = ["name", "description", "source", "path", "exclude"]
    for key in rep_keys:
      if key in node.data:
        rep[key] = node.data[key]
    if len(node.children) >= 1:
      rep["children"] = sorted(node.children.values(), key=lambda mod: mod.name())
    return representer.represent_mapping(cls.yaml_tag, rep)

  @classmethod
  def from_yaml(cls, constructor, node):
    data = constructor.construct_mapping(node, deep=True)
    children = data.pop("children", [])
    item = cls(data)
    for child in children:
      item.add_child(child)
    return item

def generate_index(root, org_name):
  url = f"{TERRAFORM_REGISTRY_BASE}/modules/{org_name}"
  r = requests.get(url, params={"limit": 100})
  data = r.json()

  for module in data.get("modules", []):
    item = root.add_child_data(module)
    id = module.get("id")
    r = requests.get(f"{TERRAFORM_REGISTRY_BASE}/modules/{id}")
    data = r.json()
    for child in data.get("submodules", []):
      item.add_child_data(child)

def render_index(index, templates_dir, docs_dir):
  env = Environment(
    keep_trailing_newline=True,
    loader=FileSystemLoader(docs_dir),
    trim_blocks=True,
    lstrip_blocks=True,
  )
  templates = env.list_templates()

  for template_file in templates:
    if not template_file.endswith(".tmpl"):
      continue
    output_file = os.path.basename(template_file.replace(".tmpl", ""))

    template = env.get_template(template_file)
    modules = [mod for mod in index.children.values() if mod.should_display()]
    rendered = template.render(modules=modules)

    with open(os.path.join(docs_dir, output_file), "w") as f:
      f.write(rendered)

def main(argv):
  parser = argparser()
  args = parser.parse_args(argv[1:])
  docs_dir = args.docs_dir
  meta_dir = os.path.join(docs_dir, "meta")

  index_file = os.path.join(meta_dir, "index.yaml")
  with open(index_file, "r+") as f:
    root = yaml.load(f, Loader=yaml.Loader)

    if not args.skip_refresh:
      generate_index(root, "terraform-google-modules")
      generate_index(root, "googlecloudplatform")

    f.seek(0)
    f.truncate()
    yaml.dump(root, f)

  render_index(root, meta_dir, docs_dir)

def argparser():
  parser = argparse.ArgumentParser(description='Generate index of blueprints')
  parser.add_argument('docs_dir', metavar='F')

  parser.add_argument('--skip-refresh', default=False, action='store_true')

  return parser

if __name__ == "__main__":
    main(sys.argv)
