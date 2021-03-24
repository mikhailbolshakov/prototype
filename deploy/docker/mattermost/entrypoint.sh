#!/bin/sh

src_plugins="./src/plugins"
trg_plugins="./plugins"
src_config="./src/config/config.json"
trg_config="./config/config.json"

echo "checking and coping prepackaged plugins..."
for entry in `ls ${src_plugins}`
do
  src_file="$src_plugins/$entry"
  trg_file="$trg_plugins/$entry"
  if [ ! -e $trg_file ]; then
    echo "$trg_file doesn't exist, copy to $trg_plugins"
    cp -r $src_file $trg_plugins
  else
    echo "$trg_file already exists, skip"
  fi
done
echo "plugins processed"

echo "checking and coping config..."
if [ ! -f $trg_config ]; then
  echo "$trg_config doesn't exist, copy"
  cp $src_config $trg_config
else
  echo "$trg_config exists, skip"
fi
echo "config processed"

_1=$(echo "$1" | awk '{ s=substr($0, 0, 1); print s; }')
if [ "$_1" = '-' ]; then
  set -- mattermost "$@"
fi

# starting mattermost
exec "$@"