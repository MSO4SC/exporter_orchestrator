#!/bin/sh

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

SESSION=$(echo $2 | sed 's/\./-/g') # Get HOST as session name (container name in this case)
NAME="slurmExp_"$SESSION

docker run -d -p 9100 --name $NAME \
          mso4sc/slurm_exporter \
            -host $2 -ssh-user $3 -ssh-password $4 -countrytz $5 -log-level=$6

status=$?
if [ $status == 0 ]; then
  PORT="$(docker ps |grep $NAME|sed 's/.*0.0.0.0://g'|sed 's/->.*//g')"
  cat > /mso4sc/targets/$SESSION.json <<- EOM
[
  {
    "targets": ["localhost:$PORT"],
    "labels": {
      "env": "canary",
      "job": "$2"
    }
  }
]
EOM
fi