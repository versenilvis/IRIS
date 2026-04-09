# iris zsh fast IPC integration

if [[ -n "$IRIS_FD" ]]; then
  _iris_send_lbuffer() {
    # -u $IRIS_FD writes to the pipe file descriptor set up by Iris wrappers
    # -N appends a null byte '\0' instead of newline (perfect for parsing)
    # -r prints raw string
    print -u $IRIS_FD -N -r -- "$LBUFFER" 2>/dev/null
  }

  autoload -Uz add-zle-hook-widget
  # Hook into ZLE so this runs absolutely every time the line buffer changes
  add-zle-hook-widget line-pre-redraw _iris_send_lbuffer
fi
