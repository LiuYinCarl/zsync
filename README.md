# zsync
目录同步工具

开发这个工具的需要来源于如下，平时办公使用到多台设备，移动硬盘，会将一些资料/文档等保存在不同的移动硬盘中，但是文档
积累越来越多，整理起来非常麻烦。Git, SVN 工具不适用这个需求，
- 一方面是硬盘大小不同，文档的重要性不同，不需要做全量同步，Git, SVN 等工具不太适合这个需求。
- 另一方面 Git, SVN 会在磁盘上添加版本目录 `.git/`, `.svn/`，属于不必要的磁盘占用。

zsync 工具的原理比较简单，用户给出源目录和目的目录，工具递归访问两个目录下的文件，计算文件的 MD5, 然后列出源目录存在但是目的目录
不存在的文件，如果添加了 `-c=true` 参数，则把这些差异文件复制到源目录下的 `zsync_temp/` 目录下，给用户自己进行整理。

这个项目包含两个程序，`zsync` 和 `zsync_coroutine`，两个程序的功能完全一致，`zsync_coroutine` 使用了协程特性，可以通过添加 `-g` 参数来指定计算 MD5 的时候的
协程数量，在文件数量的文件大小比较大的情况下，有很好的加速效果，300 个文件占用 8 G 磁盘空间的情况下使用 `-g=20` 参数的运行时间大约是 `-g=1` 参数的运行时间的 1/4。

## 编译运行例子

这是一个 Golang 项目，所以需要使用 Golang 编译器进行编译。

```bash
# 编译
go build -o zsync.exe zsync.go
# 查看帮助信息
./zsync.exe --help
Usage of zsync.exe:
  -c    copy to temp dir in destination directory
  -dst string
        absolute path of destination directory
  -src string
        absolute path of source directory
# 运行程序
./zsync.exe -src="D:\ebooks" -dst="E:ebooks"


go build -o zsync_coroutine.exe zsync_coroutine.go
./zsync_coroutine.exe --help
Usage of zsync_coroutine.exe:
  -c    copy to temp dir in destination directory, default=false
  -dst string
        absolute path of destination directory
  -g int
        coroutine count, default=10 (default 10)
  -src string
        absolute path of source directory
# 运行程序
./zsync_coroutine.exe -src="D:\ebooks" -dst="E:ebooks"
```

## 需要注意的地方

1. 如果源目录下的不同子目录有多个同名文件，那么使用 `-c` 参数的时候，同名文件会被覆盖，所以如果发现了同名文件，建议在运行 zsync 之后及时清理 `zsync_temp/` 目录下的文件，之后再运行一次 zsync。目前出于不修改原文件名的考虑不对这个问题做处理。


