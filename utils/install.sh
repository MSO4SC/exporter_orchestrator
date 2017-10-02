#!/bin/bash

# Copyright 2017 MSO4SC - javier.carnero@atos.net
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if [[ $# -eq 0 ]] ; then
    echo 'Usage: '$0' WORKDIR' 
    exit 1
fi

go install github.com/mso4sc/exporter_orchestrator

cp $GOPATH/src/github.com/mso4sc/exporter_orchestrator/config.json $1/config.json
mkdir $1/scripts
cp $GOPATH/src/github.com/mso4sc/exporter_orchestrator/scripts/* $1/scripts
chmod +x $1/scripts/*