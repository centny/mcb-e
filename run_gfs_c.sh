#!/bin/bash
cd `dirname ${0}`
export PATH=`pwd`:`dirname ${0}`:$PATH
gfs -c mcb_c.properties