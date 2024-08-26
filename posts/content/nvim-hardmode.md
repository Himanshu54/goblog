---
title: "nvim hardmode in init.lua"
date: 2023-08-05T16:02:11Z
draft: false
tags:
    - nvim
    - vim
    - lua
    - ".dotfile"
---

In my quest to learn vim, I stumbled upon a very excellent advise to disable arrow keys and force yourself to use only h,j,k,l movements. Being a former counter-strike enthusiast I know that this method works as I have forced me to learn better keybinds earlier too.

And since I am also looking into how to configure nvim properly whilest understanding what I put into my configurations, I decided not to blindly add [vim-hardmode](https://github.com/dusans/vim-hardmode) plugin.

First I though I would just copy over the keymaps directive to my newly created `init.lua` config, but the plugin had the binds in native vim directives. So to keep up with my resolve to use `lua` to configure my nvim, I had to rewrite the instruction in `lua` ( sort off ).

Thankfully [Neovim lua user guide](https://neovim.io/doc/user/lua-guide.html) is quite enought to do the conversions. Here is the result of the conversion.
```lua
-- Bindings to enable hardmode on vim, break habits.
-- https://github.com/dusans/vim-hardmode ported to lua.
-- TODO: Disable h,j,k,l in normal and visual mode ??
-- START HARDMODE
vim.g.hardmodemsg = "VIM: hard Mode [ ' :call EasyMode()' to exit]"
vim.g.hardmode_on = 0

function HardMode()
	vim.opt.backspace=""
	vim.keymap.set({'n', 'v', 'i'}, '<Left>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.keymap.set({'n', 'v', 'i'}, '<Right>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.keymap.set({'n', 'v', 'i'}, '<Up>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.keymap.set({'n', 'v', 'i'}, '<Down>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.keymap.set({'n', 'v', 'i'}, '<PageUp>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.keymap.set({'n', 'v', 'i'}, '<PageDown>', "<Esc>:echo g:hardmodemsg<CR>", { buffer=true })
	vim.g.hardmode_on = 1
	print(vim.g.hardmodemsg)
end

function EasyMode()
	vim.opt.backspace="indent,eol,start"
	vim.keymap.del({'n', 'v', 'i'}, '<Left>', { buffer=true, silent = true })
	vim.keymap.del({'n', 'v', 'i'}, '<Right>', { buffer=true, silent = true })
	vim.keymap.del({'n', 'v', 'i'}, '<Up>', { buffer=true, silent = true })
	vim.keymap.del({'n', 'v', 'i'}, '<Down>', { buffer=true, silent = true })
	vim.keymap.del({'n', 'v', 'i'}, '<PageUp>', { buffer=true, silent = true })
	vim.keymap.del({'n', 'v', 'i'}, '<PageDown>', { buffer=true, silent = true })
	vim.g.hardmode_on = 0
	print("You are weak...")
end

function ToggleMode()
	if (vim.g.hardmode_on==1)
	then
		EasyMode()
	else
		HardMode()
	end
end
vim.api.nvim_create_user_command(
  'ToggleMode',
  ToggleMode,
  { nargs=0}
  )
```

Someone who have looked into the hardmode config for vim, might find that the directive are compact comparitively due the fact that I can set for `normal`, `visual` and `insert` mode simultaneously in a single line.
