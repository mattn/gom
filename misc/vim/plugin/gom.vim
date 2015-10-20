let s:save_cpo = &cpo
set cpo&vim


function! s:setGomEnv()
  let $GOPATH = filter(split(system("gom exec env"), "\n"), "v:val =~ '^GOPATH='")[0][7:]
endfunction


command! SetGomEnv call s:setGomEnv()


let &cpo = s:save_cpo
unlet s:save_cpo
