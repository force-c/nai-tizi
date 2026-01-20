# macOS 开发环境配置文档

本目录包含 macOS 系统相关的配置和问题排查文档。

## 文档列表

### [InfluxDB CLI 问题排查](./influx-cli-troubleshooting.md)
- InfluxDB CLI 命令找不到的原因分析
- 详细的排查步骤
- PATH 配置修复方法
- 预防措施

### [终端主题配置指南](./terminal-theme-configuration.md)
- iTerm2 和 Terminal.app 使用不同主题的配置方法
- agnoster 主题配置说明
- 系统默认提示符恢复
- Powerline 字体配置
- 常见问题解答

## 适用场景

这些文档适用于以下情况：

1. **终端主题配置后命令找不到** - 参考 InfluxDB CLI 问题排查文档
2. **想在不同终端使用不同主题** - 参考终端主题配置指南
3. **Powerline 字体显示问题** - 参考终端主题配置指南
4. **恢复 macOS 默认终端样式** - 参考终端主题配置指南

## 相关工具

- **Homebrew** - macOS 包管理器
- **Oh My Zsh** - Zsh 配置框架
- **iTerm2** - 增强型终端应用
- **Terminal.app** - macOS 系统默认终端

## 快速参考

### 常用命令

```bash
# 重新加载 Zsh 配置
source ~/.zshrc

# 查看当前 PATH
echo $PATH

# 查看当前终端类型
echo $TERM_PROGRAM

# 查看当前主题
echo $ZSH_THEME

# 备份配置文件
cp ~/.zshrc ~/.zshrc.backup.$(date +%Y%m%d_%H%M%S)
```

### 配置文件位置

- `~/.zshrc` - Zsh 主配置文件
- `~/.zprofile` - Zsh 登录配置
- `~/.oh-my-zsh/` - Oh My Zsh 目录
- `~/Library/Fonts/` - 用户字体目录
