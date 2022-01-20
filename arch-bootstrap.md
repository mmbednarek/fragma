# Steps for bootstrapping arch linux image

1. Download package from archlinux repository
2. Setup resolve server
```
nameserver 8.8.8.8
nameserver 8.8.4.4
```
3. Setup pacman mirrors
4. Run pacman-key --init
5. Run pacman-key --populate archlinux
6. Run pacman -Syyu