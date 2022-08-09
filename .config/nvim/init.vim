let skip_defaults_vim=1
set viminfo=""
set viminfofile=NONE
set showcmd
set mouse=a
set splitright
set splitbelow
set ignorecase
set smartcase
set nocursorcolumn
set pumheight=20
set tabstop=4
set shiftwidth=4
set shortmess+=c
set maxmempattern=400000
set completeopt=menu,menuone,longest 
set number
set encoding=utf-8
set autoread
set autoindent
set incsearch
set hlsearch

autocmd TermOpen * setlocal nonumber norelativenumber

set undofile
set undodir=/tmp/vim/undofile

syntax enable

" install
call plug#begin()
Plug 'preservim/nerdtree'
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
Plug 'lewis6991/gitsigns.nvim'
call plug#end()

lua require('gitsigns').setup()

" keybind to open

nnoremap <C-n> :NERDTree<CR>

" vim-go

" write before running other vim-go commands
set autowrite

" some hotkeys

autocmd FileType go nmap <leader>b  <Plug>(go-build)
autocmd FileType go nmap <leader>r  <Plug>(go-run)
autocmd FileType go nmap <leader>t  <Plug>(go-test)

" open auto compelete on .
au filetype go inoremap <buffer> . .<C-x><C-o>

" setup syntax highlighting
let g:go_doc_popup_window = 1
let g:go_auto_sameids = 1
let g:go_highlight_array_whitespace_error = 1
let g:go_highlight_chan_whitespace_error = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_space_tab_error = 1
let g:go_highlight_trailing_whitespace_error = 0
let g:go_highlight_operators = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_parameters = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_build_constraints = 1
let g:go_highlight_generate_tags = 1
let g:go_highlight_string_spellcheck = 1
let g:go_highlight_format_strings = 1
let g:go_highlight_variable_declarations = 1
let g:go_highlight_variable_assignments = 1
let g:go_fmt_experimental = 1
let g:go_metalinter_autosave=1
let g:go_metalinter_autosave_enabled=['golint', 'govet']