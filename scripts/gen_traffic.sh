#!/bin/bash

while true; do
  curl http://localhost/foo
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/bar
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/foo/1
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/bar/1
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/foo/2
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/
  sleep .$[ ( $RANDOM % 10 ) + 1 ]s
  curl http://localhost/
done
