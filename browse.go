package xmux

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var defaultTemplates = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Name }} - 文件浏览器</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://cdn.jsdelivr.net/npm/font-awesome@4.7.0/css/font-awesome.min.css" rel="stylesheet">
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#3b82f6',
                        secondary: '#64748b',
                        neutral: '#f1f5f9',
                        dark: '#0f172a',
                    },
                    fontFamily: {
                        sans: ['Inter', 'system-ui', 'sans-serif'],
                    },
                }
            }
        }
    </script>
    <style type="text/tailwindcss">
        @layer utilities {
            .content-auto {
                content-visibility: auto;
            }
            .file-grid {
                display: grid;
                grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
                gap: 1.5rem;
            }
            .transition-all-300 {
                transition: all 300ms ease-in-out;
            }
        }
    </style>
</head>
<body class="bg-gray-50 text-gray-800 min-h-screen">
    <!-- 顶部导航栏 -->
    <header class="bg-white shadow-sm sticky top-0 z-10">
        <div class="container mx-auto px-4 py-4 flex flex-col md:flex-row justify-between items-start md:items-center">
            <div class="flex items-center mb-4 md:mb-0">
                <i class="fa fa-folder-open text-primary text-2xl mr-3"></i>
                <h1 class="text-xl font-bold">文件浏览器</h1>
            </div>
            
            <!-- 面包屑导航 - 修复了loop函数未定义的问题 -->
            <nav class="w-full md:w-auto text-sm text-gray-600">
                {{ $len := len .Breadcrumbs }}
                {{ range $index, $crumb := .Breadcrumbs }}
                    {{ if eq $index (sub $len 1) }}
                        <span class="text-primary font-medium">{{ .Text }}</span>
                    {{ else }}
                        <a href="?path={{ .Link }}" class="hover:text-primary transition-all-300">{{ .Text }}</a>
                        <span class="mx-2">/</span>
                    {{ end }}
                {{ end }}
            </nav>
        </div>
    </header>

    <main class="container mx-auto px-4 py-6">
        <!-- 工具栏 -->
        <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div class="flex items-center space-x-4">
                    {{ if .CanGoUp }}
                        <a href="{{ .Prev }}" class="inline-flex items-center px-3 py-1.5 bg-primary text-white rounded-md hover:bg-primary/90 transition-all-300">
                            <i class="fa fa-arrow-up mr-2"></i> 上一级
                        </a>
                    {{ end }}
                    
                    <div class="text-gray-600">
                        <span class="mr-4"><i class="fa fa-folder text-yellow-500 mr-1"></i> {{ .NumDirs }} 个目录</span>
                        <span><i class="fa fa-file text-gray-400 mr-1"></i> {{ .NumFiles }} 个文件</span>
                    </div>
                </div>
                
                <div class="flex items-center space-x-3 w-full md:w-auto">
                    <div class="text-sm text-gray-500">
                        总大小: <span class="font-medium">{{ .HumanTotalFileSize }}</span>
                    </div>
                    
                    <!-- 布局切换 -->
                    <div class="border rounded-md p-1 flex">
                        <a href="?path={{ .Path }}&layout=list&sort={{ .Sort }}&order={{ .Order }}" 
                           class="px-2 py-1 rounded {{ if eq .Layout "list" }}bg-primary text-white{{ else }}text-gray-600 hover:bg-gray-100{{ end }}">
                            <i class="fa fa-list"></i>
                        </a>
                        <a href="?path={{ .Path }}&layout=grid&sort={{ .Sort }}&order={{ .Order }}" 
                           class="px-2 py-1 rounded {{ if eq .Layout "grid" }}bg-primary text-white{{ else }}text-gray-600 hover:bg-gray-100{{ end }}">
                            <i class="fa fa-th-large"></i>
                        </a>
                    </div>
                </div>
            </div>
            
            <!-- 排序选项 -->
            <div class="mt-4 pt-4 border-t border-gray-100 flex items-center text-sm">
                <span class="text-gray-500 mr-3">排序方式:</span>
                <div class="flex space-x-4">
                    <a href="?path={{ .Path }}&layout={{ .Layout }}&sort=name&order={{ if and (eq .Sort "name") (eq .Order "asc") }}desc{{ else }}asc{{ end }}"
                       class="flex items-center {{ if eq .Sort "name" }}text-primary{{ else }}text-gray-600 hover:text-primary{{ end }}">
                        名称
                        {{ if eq .Sort "name" }}
                            <i class="fa ml-1 {{ if eq .Order "asc" }}fa-sort-asc{{ else }}fa-sort-desc{{ end }}"></i>
                        {{ end }}
                    </a>
                    <a href="?path={{ .Path }}&layout={{ .Layout }}&sort=size&order={{ if and (eq .Sort "size") (eq .Order "asc") }}desc{{ else }}asc{{ end }}"
                       class="flex items-center {{ if eq .Sort "size" }}text-primary{{ else }}text-gray-600 hover:text-primary{{ end }}">
                        大小
                        {{ if eq .Sort "size" }}
                            <i class="fa ml-1 {{ if eq .Order "asc" }}fa-sort-asc{{ else }}fa-sort-desc{{ end }}"></i>
                        {{ end }}
                    </a>
                    <a href="?path={{ .Path }}&layout={{ .Layout }}&sort=modtime&order={{ if and (eq .Sort "modtime") (eq .Order "asc") }}desc{{ else }}asc{{ end }}"
                       class="flex items-center {{ if eq .Sort "modtime" }}text-primary{{ else }}text-gray-600 hover:text-primary{{ end }}">
                        修改时间
                        {{ if eq .Sort "modtime" }}
                            <i class="fa ml-1 {{ if eq .Order "asc" }}fa-sort-asc{{ else }}fa-sort-desc{{ end }}"></i>
                        {{ end }}
                    </a>
                </div>
            </div>
        </div>
        
        <!-- 文件列表 - 列表视图 -->
        {{ if eq .Layout "list" }}
            <div class="bg-white rounded-lg shadow-sm overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">大小</th>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">修改时间</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        {{ range .Items }}
                        <tr class="hover:bg-gray-50 transition-all-300">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center">
                                        {{ if .IsDir }}
                                            <i class="fa fa-folder text-2xl text-yellow-500"></i>
                                        {{ else if .IsSymlink }}
                                            <i class="fa fa-link text-2xl text-blue-500"></i>
                                        {{ else if HasExt .Ext "jpg" "jpeg" "png" "gif" "bmp" }}
                                            <i class="fa fa-file-image-o text-2xl text-green-500"></i>
                                        {{ else if HasExt .Ext "pdf" }}
                                            <i class="fa fa-file-pdf-o text-2xl text-red-500"></i>
                                        {{ else if HasExt .Ext "doc" "docx" }}
                                            <i class="fa fa-file-word-o text-2xl text-blue-700"></i>
                                        {{ else if HasExt .Ext "xls" "xlsx" }}
                                            <i class="fa fa-file-excel-o text-2xl text-green-600"></i>
                                        {{ else if HasExt .Ext "zip" "rar" "tar" "gz" }}
                                            <i class="fa fa-file-archive-o text-2xl text-purple-500"></i>
                                        {{ else if HasExt .Ext "txt" "md" "html" "css" "js" "go" }}
                                            <i class="fa fa-file-text-o text-2xl text-gray-600"></i>
                                        {{ else }}
                                            <i class="fa fa-file-o text-2xl text-gray-400"></i>
                                        {{ end }}
                                    </div>
                                    <div class="ml-4">
                                        <div class="text-sm font-medium text-gray-900">
                                            <a href="{{ .URL }}" class="hover:text-primary transition-all-300">
                                                {{ .Name }}
                                                {{ if .IsSymlink }}
                                                    <span class="text-xs text-gray-500 ml-2">(链接到: {{ .SymlinkPath }})</span>
                                                {{ end }}
                                            </a>
                                        </div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ if .IsDir }}-{{ else }}{{ .HumanSize }}{{ end }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ .HumanModTime  }}
                            </td>
                        </tr>
                        {{ else }}
                        <tr>
                            <td colspan="3" class="px-6 py-8 text-center text-gray-500">
                                <i class="fa fa-folder-open-o text-4xl mb-2 opacity-30"></i>
                                <p>当前目录为空</p>
                            </td>
                        </tr>
                        {{ end }}
                    </tbody>
                </table>
            </div>
        {{ end }}
        
        <!-- 文件列表 - 网格视图 -->
        {{ if eq .Layout "grid" }}
            <div class="file-grid">
                {{ range .Items }}
                <div class="bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-all-300">
                    <a href="{{ .URL }}" class="block h-full">
                        <div class="p-6 flex flex-col items-center justify-center bg-gray-50">
                            {{ if .IsDir }}
                                <i class="fa fa-folder text-5xl text-yellow-500"></i>
                            {{ else if .IsSymlink }}
                                <i class="fa fa-link text-5xl text-blue-500"></i>
                            {{ else if HasExt .Ext "jpg" "jpeg" "png" "gif" "bmp" }}
                                <i class="fa fa-file-image-o text-5xl text-green-500"></i>
                            {{ else if HasExt .Ext "pdf" }}
                                <i class="fa fa-file-pdf-o text-5xl text-red-500"></i>
                            {{ else if HasExt .Ext "doc" "docx" }}
                                <i class="fa fa-file-word-o text-5xl text-blue-700"></i>
                            {{ else if HasExt .Ext "xls" "xlsx" }}
                                <i class="fa fa-file-excel-o text-5xl text-green-600"></i>
                            {{ else if HasExt .Ext "zip" "rar" "tar" "gz" }}
                                <i class="fa fa-file-archive-o text-5xl text-purple-500"></i>
                            {{ else if HasExt .Ext "txt" "md" "html" "css" "js" "go" }}
                                <i class="fa fa-file-text-o text-5xl text-gray-600"></i>
                            {{ else }}
                                <i class="fa fa-file-o text-5xl text-gray-400"></i>
                            {{ end }}
                        </div>
                        <div class="p-3">
                            <div class="text-sm font-medium text-gray-900 truncate" title="{{ .Name }}">
                                {{ .Name }}
                            </div>
                            <div class="mt-1 text-xs text-gray-500">
                                {{ if .IsDir }}
                                    目录
                                {{ else }}
                                    {{ .HumanSize }}
                                {{ end }}
                            </div>
                            <div class="mt-1 text-xs text-gray-500">
                                {{ .HumanModTime  }}
                            </div>
                        </div>
                    </a>
                </div>
                {{ else }}
                <div class="col-span-full text-center py-12 text-gray-500">
                    <i class="fa fa-folder-open-o text-5xl mb-3 opacity-30"></i>
                    <p>当前目录为空</p>
                </div>
                {{ end }}
            </div>
        {{ end }}
    </main>

    <footer class="bg-white border-t mt-8 py-6">
        <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
            <p>文件浏览器 &copy; 2023</p>
        </div>
    </footer>
</body>
</html>
`

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".zip":
		return "application/zip"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream" // 二进制流，通用类型
	}
}

func FileBrowse(pattern, rootDir string, listDir bool, download bool) *RouteGroup {
	g := NewRouteGroup().BindResponse(nil)

	g.Get(pattern+"{all:filename}", func(w http.ResponseWriter, r *http.Request) {
		path := Var(r)["filename"]
		fullPath := filepath.Join(rootDir, path)
		prefix := r.URL.Path
		// 安全检查：防止目录遍历攻击
		if !isSubdirectory(rootDir, fullPath) {
			http.Error(w, "访问被拒绝", http.StatusForbidden)
			return
		}

		// 检查路径是否存在
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			http.Error(w, "路径不存在", http.StatusNotFound)
			return
		}

		// 如果是文件，直接下载
		if !fileInfo.IsDir() {
			file, _ := os.Open(fullPath)
			defer file.Close()
			if download {
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
			}

			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
			http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
			return
		}
		if !listDir {
			w.WriteHeader(404)
			return
		}
		// 获取布局和排序参数
		layout := r.URL.Query().Get("layout")
		if layout != "grid" {
			layout = "list"
		}

		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			sortBy = "name"
		}

		order := r.URL.Query().Get("order")
		if order == "" {
			order = "asc"
		}

		// 生成TemplateData（实际实现应使用之前的PopulateTemplateData函数）
		data, err := PopulateTemplateData(prefix, rootDir, fullPath, 0, 0)
		if err != nil {
			http.Error(w, "无法读取目录内容", http.StatusInternalServerError)
			return
		}
		urlPath := "/" + strings.ReplaceAll(prefix, "\\", "/")
		data.Prev = strings.ReplaceAll(filepath.Join(prefix, "../"), "\\", "/")
		if len(data.Prev) < len(pattern) {
			data.Prev = urlPath
		}
		data.Layout = layout
		data.Sort = sortBy
		data.Order = order
		// 渲染模板
		funcMap := template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			"HasExt": HasExt,
		}

		tmpl := template.Must(template.New("browse.html").Funcs(funcMap).Parse(defaultTemplates))

		err = tmpl.Execute(w, data)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	})

	return g
}

// 定义结构体（与之前保持一致）
type TemplateData struct {
	CSP        string
	EnableCsp  bool
	RespHeader struct {
		Set func(string, string)
	}
	Name        string
	Path        string
	Breadcrumbs []struct {
		Link string
		Text string
	}
	NumDirs            int
	NumFiles           int
	HumanTotalFileSize string
	Limit              int
	Offset             int
	Layout             string // "list" or "grid"
	Sort               string
	Order              string
	CanGoUp            bool
	Items              []FileItem
	Prev               string
}

type FileItem struct {
	Name         string
	URL          string
	IsDir        bool
	IsSymlink    bool
	SymlinkPath  string
	Size         int64
	HumanSize    string
	ModTime      time.Time
	HumanModTime string
	HasExt       func(...string) bool
	Tpl          struct {
		Layout string
	}
	Ext string
}

// 检查子目录关系，防止目录遍历
func isSubdirectory(parent, child string) bool {
	absParent, err := filepath.Abs(parent)
	if err != nil {
		return false
	}

	absChild, err := filepath.Abs(child)
	if err != nil {
		return false
	}

	return hasPathPrefix(absChild, absParent)
}

// 判断路径是否以指定前缀开头（替代 deprecated 的 filepath.HasPrefix）
func hasPathPrefix(path, prefix string) bool {
	// 标准化路径（处理 . 和 .. 等特殊符号）
	cleanPath := filepath.Clean(path)
	cleanPrefix := filepath.Clean(prefix)

	// 在 Windows 系统上，路径大小写不敏感，需要转为统一大小写判断
	if filepath.Separator == '\\' { // 判断是否为 Windows 系统
		cleanPath = strings.ToLower(cleanPath)
		cleanPrefix = strings.ToLower(cleanPrefix)
	}

	// 判断前缀，同时确保前缀是一个完整的路径段
	return strings.HasPrefix(cleanPath, cleanPrefix) &&
		(len(cleanPath) == len(cleanPrefix) ||
			strings.Index(cleanPath[len(cleanPrefix):], string(filepath.Separator)) == 0)
}

// PopulateTemplateData 扫描目录并填充TemplateData结构体
func PopulateTemplateData(pattern, rootDir, currentDir string, limit, offset int) (*TemplateData, error) {
	// 转换为绝对路径
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("获取根目录绝对路径失败: %w", err)
	}

	absCurrent, err := filepath.Abs(currentDir)
	if err != nil {
		return nil, fmt.Errorf("获取当前目录绝对路径失败: %w", err)
	}

	// 安全检查：确保当前目录在根目录之下
	if !strings.HasPrefix(absCurrent, absRoot) {
		return nil, fmt.Errorf("访问被拒绝：当前目录不在根目录范围内")
	}

	// 初始化TemplateData基础字段
	td := &TemplateData{
		EnableCsp: true,
		CSP:       "default-src 'self'; script-src 'nonce-%s'; style-src 'nonce-%s' https://cdn.tailwindcss.com; font-src https://cdn.jsdelivr.net;",
		Path:      currentDir,
		Name:      filepath.Base(currentDir),
		Limit:     limit,
		Offset:    offset,
		Layout:    "list",                // 默认列表布局
		Sort:      "name",                // 默认按名称排序
		Order:     "asc",                 // 默认升序
		CanGoUp:   absCurrent != absRoot, // 是否可以向上导航
	}

	// 初始化响应头设置函数
	td.RespHeader.Set = func(key, value string) {
		// 实际使用时会在HTTP处理函数中关联到真实的ResponseWriter
	}

	// 生成面包屑导航
	td.Breadcrumbs = generateBreadcrumbs(absRoot, absCurrent)

	// 读取目录内容
	entries, err := os.ReadDir(absCurrent)
	if err != nil {
		return nil, fmt.Errorf("读取目录内容失败: %w", err)
	}

	// 处理目录项并填充FileItem
	var totalFileSize int64
	var fileItems []FileItem

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的条目
		}

		// 计算总文件大小（仅统计文件）
		if !entry.IsDir() {
			totalFileSize += info.Size()
		}

		// 创建文件项
		item := FileItem{
			Name:      info.Name(),
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			HumanSize: humanReadableSize(info.Size()),
			ModTime:   info.ModTime(),
			Ext:       strings.ToLower(filepath.Ext(info.Name())),
			Tpl:       struct{ Layout string }{Layout: td.Layout},
		}

		// 处理符号链接
		if info.Mode()&os.ModeSymlink != 0 {
			item.IsSymlink = true
			target, err := os.Readlink(filepath.Join(absCurrent, info.Name()))
			if err == nil {
				item.SymlinkPath = target
			}
		}

		// 生成相对URL
		item.URL = filepath.ToSlash(filepath.Join(pattern, info.Name()))
		// if err == nil {
		// 	// 转换为URL路径格式（使用斜杠）
		// item.URL = filepath.ToSlash(relPath)
		// }

		// 初始化FileItem的方法
		item.initHumanModTime()
		item.initHasExt()

		fileItems = append(fileItems, item)
	}

	// 排序文件项
	td.sortItems(fileItems)

	// 应用分页
	td.Items = applyPagination(fileItems, limit, offset)

	// 填充统计信息
	td.NumFiles = countFiles(td.Items)
	td.NumDirs = len(td.Items) - td.NumFiles
	td.HumanTotalFileSize = humanReadableSize(totalFileSize)

	return td, nil
}

// 生成面包屑导航
func generateBreadcrumbs(root, current string) []struct{ Link, Text string } {
	breadcrumbs := []struct{ Link, Text string }{
		{Link: "", Text: filepath.Base(root)},
	}

	if root == current {
		return breadcrumbs
	}

	relPath, err := filepath.Rel(root, current)
	if err != nil {
		return breadcrumbs
	}

	parts := strings.Split(relPath, string(os.PathSeparator))
	currentPath := root

	for _, part := range parts {
		currentPath = filepath.Join(currentPath, part)
		relToRoot, err := filepath.Rel(root, currentPath)
		if err != nil {
			continue
		}

		breadcrumbs = append(breadcrumbs, struct{ Link, Text string }{
			Link: filepath.ToSlash(relToRoot),
			Text: part,
		})
	}

	return breadcrumbs
}

// 人类可读的文件大小
func humanReadableSize(size int64) string {
	switch {
	case size >= 1<<40: // 1 TB
		return fmt.Sprintf("%.2f TB", float64(size)/(1<<40))
	case size >= 1<<30: // 1 GB
		return fmt.Sprintf("%.2f GB", float64(size)/(1<<30))
	case size >= 1<<20: // 1 MB
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10: // 1 KB
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// 统计文件数量（排除目录）
func countFiles(items []FileItem) int {
	count := 0
	for _, item := range items {
		if !item.IsDir {
			count++
		}
	}
	return count
}

// 应用分页
func applyPagination(items []FileItem, limit, offset int) []FileItem {
	if limit <= 0 || offset < 0 {
		return items
	}

	start := offset
	end := offset + limit

	if start >= len(items) {
		return []FileItem{}
	}

	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

// 排序文件项
func (td *TemplateData) sortItems(items []FileItem) {
	sort.Slice(items, func(i, j int) bool {
		// 目录排在文件前面
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}

		// 根据指定字段排序
		switch td.Sort {
		case "size":
			if td.Order == "desc" {
				return items[i].Size > items[j].Size
			}
			return items[i].Size < items[j].Size
		case "modtime":
			if td.Order == "desc" {
				return items[i].ModTime.After(items[j].ModTime)
			}
			return items[i].ModTime.Before(items[j].ModTime)
		case "name", "":
			fallthrough
		default:
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}
	})

	// 名称降序排序
	if td.Sort == "name" && td.Order == "desc" {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
	}
}

// 初始化HumanModTime方法
func (fi *FileItem) initHumanModTime() {
	fi.HumanModTime = fi.ModTime.Format("2006-01-02 15:04:05")
}

// 初始化HasExt方法
func (fi *FileItem) initHasExt() {
	fi.HasExt = func(exts ...string) bool {
		if fi.Ext == "" {
			return false
		}
		ext := strings.TrimPrefix(fi.Ext, ".") // 移除点号
		for _, e := range exts {
			if strings.TrimPrefix(strings.ToLower(e), ".") == strings.ToLower(ext) {
				return true
			}
		}
		return false
	}
}

func HasExt(ext string, exts ...string) bool {
	if ext == "" {
		return false
	}

	for _, e := range exts {
		if strings.TrimPrefix(strings.ToLower(e), ".") == strings.ToLower(ext) {
			return true
		}
	}
	return false
}
