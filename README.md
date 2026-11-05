# my-tree

Personal website — [emilio.eguia.me](https://emilio.eguia.me)

Built with [Hugo](https://gohugo.io) · Pure CSS · [Cloudflare Pages](https://pages.cloudflare.com)

**No JavaScript. Ever.**

---

## Stack

| Layer | Choice |
|-------|--------|
| SSG | Hugo 0.147.9 |
| Styling | Pure CSS — `@keyframes`, custom properties, no frameworks |
| Hosting | Cloudflare Pages |
| CI/CD | GitLab CI → Cloudflare Pages |
| Mirror | GitLab (primary) → GitHub |

---

## Project layout

```
my-tree/
├── templates/           # Under-construction page designs
│   ├── base/            # 10 originals — no images
│   ├── pro/             # + picsum.photos accent image (CDN, zero binaries)
│   └── meme/            # + inline SVG corner meme (~200–600 bytes)
├── errors/              # Error pages (404)
└── (hugo site — Phase 1)
```

---

## Local preview

Requirements: any browser. No build step needed for templates.

```sh
# Clone
git clone git@gitlab.com:jemo121/my-tree.git
cd my-tree

# Open any template directly
xdg-open templates/base/warm.html
```

Hugo site (Phase 1, not started):

```sh
hugo server -D
```

---

## License

[AGPL-3.0-or-later](LICENSE)
