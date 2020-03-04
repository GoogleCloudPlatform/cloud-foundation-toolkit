#!/bin/bash
kpt cfg set . folder-number $(kubectl describe -f folder.yaml | grep Name:\ *folders\/ | sed "s/.*folders\///") --set-by "set-folder-number"
