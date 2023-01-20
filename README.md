# zsync
目录同步工具

开发这个工具的需要来源于如下，平时办公使用到多台设备，移动硬盘，会将一些资料/文档等保存在不同的移动硬盘中，但是文档
积累越来越多，整理起来非常麻烦。Git, SVN 工具不适用这个需求，
- 一方面是硬盘大小不同，文档的重要性不同，不需要做全量同步，Git, SVN 等工具不太适合这个需求。
- 另一方面 Git, SVN 会在磁盘上添加版本目录 `.git/`, `.svn/`，属于不必要的磁盘占用。

zsync 工具的原理比较简单，用户给出源目录和目的目录，工具递归访问两个目录下的文件，计算文件的 MD5, 然后列出源目录存在但是目的目录
不存在的文件，如果添加了 `-c=true` 参数，则把这些差异文件复制到源目录下的 `zsync_temp/` 目录下，给用户自己进行整理。

## 编译运行例子

```bash
go build -o zsync.exe zsync.go
main.exe -src="D:\ebooks" -dst="E:ebooks"

go build -o zsync_coroutine.exe zsync_coroutine.go
main2.exe -src="D:\ebooks" -dst="E:ebooks"
```


