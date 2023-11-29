#!/bin/sh
# <swiftbar.hideAbout>true</swiftbar.hideAbout>
# <swiftbar.hideRunInTerminal>true</swiftbar.hideRunInTerminal>
# <swiftbar.hideLastUpdated>true</swiftbar.hideLastUpdated>
# <swiftbar.hideDisablePlugin>true</swiftbar.hideDisablePlugin>
# <swiftbar.hideSwiftBar>true</swiftbar.hideSwiftBar>
# <swiftbar.schedule>* * * * *</swiftbar.schedule>

# optionally configure tt
# export TT_HOME_DIR="${HOME}/.config/tt"

# path to tt binary, adjust to your system if needed
export TT_BIN="${HOME}/go/bin/tt"

# print short status to show in menu bar
"${TT_BIN}" status --short
# menu options
echo "---"
echo "Resume | bash='${TT_BIN}' param0='start' param1='--resume' terminal=false"
# Always start a static project if you don't care too much,
echo "Start | bash='${TT_BIN}' param0='start' param1='my-project' terminal=false"
# alternatively you could setup multiple 'Start X' commands.
# echo "Start Foo | bash='${TT_BIN}' param0='start' param1='foo' terminal=false"
# echo "Start Foo Bar | bash='${TT_BIN}' param0='start' param1='foo' param2='bar' terminal=false"
echo "Stop | bash='${TT_BIN}' param0='stop' terminal=false"
