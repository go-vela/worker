#!/bin/sh

DLV=$(which dlv-dap)

if [ -x "${DLV}" ]; then
	PATH="/usr/local/go/bin/:$PATH"
fi

if [ "$DEBUG_AS_ROOT" = "true" ]; then
	echo Run as Root
	# The parameter -C 4 keeps the descriptor 3 opened
	# The parameter -E is needed to ensure that environment
	# is being passed along
	#
	# Both of these also require the addition of the following
	# via visudo in the sudoers file:
	# <user> ALL=(root)NOPASSWD:SETENV:/home/<user>/go/bin/dlv-dap
	# Defaults:<user> closefrom_override
	exec sudo -E -C 4 "$DLV" --only-same-user=false "$@"
else
	echo Run as User
	exec "$DLV" "$@"
fi

