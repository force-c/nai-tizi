# macOS 终端主题配置指南

## 概述

本文档说明如何配置不同终端应用使用不同的 Zsh 主题：
- **iTerm2**：使用 agnoster 主题（Powerline 风格）
- **Terminal.app**：使用系统默认提示符

## 配置步骤

### 1. 备份配置文件

```bash
cp ~/.zshrc ~/.zshrc.backup.$(date +%Y%m%d_%H%M%S)
```

### 2. 编辑 .zshrc 配置

找到 `ZSH_THEME` 配置行（通常在文件开头），修改为：

```bash
# 注释掉原有的主题配置
# ZSH_THEME="robbyrussell"
# ZSH_THEME="agnoster"

# 根据终端类型选择主题
if [[ "$TERM_PROGRAM" == "iTerm.app" ]]; then
    ZSH_THEME="agnoster"      # iTerm2 使用 agnoster 主题
else
    ZSH_THEME=""              # 其他终端不使用主题，保持系统默认
    # 恢复系统默认提示符
    PROMPT="%n@%m %1~ %# "
fi
```

### 3. 应用配置

```bash
source ~/.zshrc
```

## 主题说明

### agnoster 主题（iTerm2）

**特点：**
- 使用 Powerline 特殊字符（箭头分隔符）
- 显示用户名、主机名、路径、Git 状态
- 需要支持 Powerline 的字体

**显示效果：**
```
 guoc@myhub  ~/dev/code_go/src  main ±✚
```

**字体要求：**
iTerm2 需要配置 Powerline 字体：
1. 打开 iTerm2 → Preferences (`⌘ + ,`)
2. Profiles → Text → Font
3. 选择带 "Powerline" 的字体，如：
   - Meslo LG S for Powerline
   - Source Code Pro for Powerline
   - Droid Sans Mono for Powerline

### 系统默认提示符（Terminal.app）

**特点：**
- macOS 原生样式
- 显示用户名@主机名、当前目录、提示符
- 不需要特殊字体

**显示效果：**
```
guoc@myhub ~ %
```

## 验证配置

### 在 iTerm2 中

```bash
echo $TERM_PROGRAM
# 输出：iTerm.app

echo $ZSH_THEME
# 输出：agnoster
```

### 在 Terminal.app 中

```bash
echo $TERM_PROGRAM
# 输出：Apple_Terminal

echo $ZSH_THEME
# 输出：（空）
```

## 常见问题

### Q1: iTerm2 中显示方块或问号

**原因：** 字体不支持 Powerline 特殊字符

**解决：** 在 iTerm2 中设置 Powerline 字体（见上文"字体要求"）

### Q2: Terminal.app 中想使用其他简单主题

可以将配置修改为：

```bash
if [[ "$TERM_PROGRAM" == "iTerm.app" ]]; then
    ZSH_THEME="agnoster"
else
    ZSH_THEME="robbyrussell"  # 或其他不需要特殊字体的主题
fi
```

推荐的简单主题：
- `robbyrussell` - Oh My Zsh 默认主题
- `simple` - 极简主题
- `clean` - 清爽主题

### Q3: 如何完全恢复到未安装 Oh My Zsh 之前

```bash
# 如果有备份文件
cp ~/.zshrc.pre-oh-my-zsh ~/.zshrc

# 或者卸载 Oh My Zsh
uninstall_oh_my_zsh
```

## 相关资源

- [Oh My Zsh 主题列表](https://github.com/ohmyzsh/ohmyzsh/wiki/Themes)
- [Powerline 字体](https://github.com/powerline/fonts)
- [agnoster 主题说明](https://github.com/agnoster/agnoster-zsh-theme)

## 配置文件位置

- `~/.zshrc` - Zsh 主配置文件
- `~/.zprofile` - Zsh 登录配置
- `~/.oh-my-zsh/` - Oh My Zsh 安装目录
- `~/.oh-my-zsh/themes/` - 主题目录
