# iris zsh fast IPC integration

if [[ -n "$IRIS_FD" ]]; then
  _iris_send_lbuffer() {
    # -u $IRIS_FD writes to the pipe file descriptor set up by Iris wrappers
    # -N appends a null byte '\0' instead of newline (perfect for parsing)
    # -r prints raw string
    print -u $IRIS_FD -N -r -- "$LBUFFER" 2>/dev/null
  }

  _iris_precmd() {
    print -u $IRIS_FD -N -r -- "IRIS_CMD_STOP" 2>/dev/null
  }

  _iris_preexec() {
    print -u $IRIS_FD -N -r -- "IRIS_CMD_START" 2>/dev/null
  }

  autoload -Uz add-zle-hook-widget
  autoload -Uz add-zsh-hook
  
  add-zle-hook-widget line-pre-redraw _iris_send_lbuffer
  add-zsh-hook precmd _iris_precmd
  add-zsh-hook preexec _iris_preexec
fi
