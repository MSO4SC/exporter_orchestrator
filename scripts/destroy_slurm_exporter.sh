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

COMMAND="/home/ubuntu/Projects/go/bin/slurm_exporter -listen-address $1 -host $2 -ssh-user $3 -ssh-password $4 -countrytz $5 -log.level=$6"
SESSION=$(echo $2 | sed 's/\./-/g')

rm /opt/prometheus/core/targets/$SESSION.json

tmux send-keys -t $SESSION:0 'C-c'
tmux kill-session -t $SESSION
