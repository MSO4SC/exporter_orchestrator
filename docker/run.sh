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

if [[ $# < 1 ]] ; then
    echo 'Usage: '$0' -monitor-host=<HOST:MPORT> [-log-level=<LOGLEVEL>]' 
    exit 1
fi

ARGS=$1
if [[ $# > 1 ]] ; then
	ARGS=$ARGS' '$2
fi
if [[ $# > 2 ]] ; then
	ARGS=$ARGS' '$3
fi

#### docker run --rm -v /lib64:/lib64 -v /usr:/usr -v /lib:/lib -v /var/run/docker.sock:/var/run/docker.sock alpine docker --version

docker run --rm -d -p 8079:8079 \
	-v /lib64:/lib64 -v /usr:/usr -v /lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu -v /var/run/docker.sock:/var/run/docker.sock \
	mso4sc/exporter_orchestrator $ARGS
