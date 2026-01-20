# InfluxDB CLI 命令找不到问题排查与修复

## 问题描述

在配置终端主题（Oh My Zsh）后，influx CLI 命令提示找不到：
```bash
influx not found
```

## 问题原因

InfluxDB@1 是 Homebrew 的 "keg-only" 包，不会自动添加到系统 PATH。在修改 `.zshrc` 配置终端主题时，可能覆盖或丢失了之前手动添加的 PATH 配置。

## 排查步骤

### 1. 确认 influx 是否已安装

```bash
# 检查 Homebrew 安装路径
brew --prefix influxdb@1

# 预期输出：/opt/homebrew/opt/influxdb@1
```

### 2. 检查文件是否存在

```bash
# 查看 influx 可执行文件
ls -la /opt/homebrew/opt/influxdb@1/bin/

# 应该能看到 influx、influxd 等文件
```

### 3. 检查当前 PATH

```bash
# 查看当前 PATH 环境变量
echo $PATH

# 检查是否包含 /opt/homebrew/opt/influxdb@1/bin
```

### 4. 检查 .zshrc 配置

```bash
# 查看是否有 influxdb 相关的 PATH 配置
grep -i influx ~/.zshrc
```

## 修复步骤

### 方案：添加 InfluxDB PATH 到 .zshrc

1. **备份当前配置**
```bash
cp ~/.zshrc ~/.zshrc.backup.$(date +%Y%m%d_%H%M%S)
```

2. **添加 InfluxDB PATH**
```bash
# 在 .zshrc 末尾添加
echo '' >> ~/.zshrc
echo '# InfluxDB PATH (added for influxdb@1)' >> ~/.zshrc
echo 'export PATH="/opt/homebrew/opt/influxdb@1/bin:$PATH"' >> ~/.zshrc
```

3. **重新加载配置**
```bash
source ~/.zshrc
```

4. **验证修复**
```bash
# 检查命令是否可用
type influx

# 预期输出：influx is /opt/homebrew/opt/influxdb@1/bin/influx

# 检查版本
influx --version

# 预期输出：InfluxDB shell version: x.x.x
```

## 预防措施

在修改 `.zshrc` 或安装终端主题时：
1. 先备份当前配置：`cp ~/.zshrc ~/.zshrc.backup`
2. 记录自定义的 PATH 配置
3. 修改后验证关键命令是否仍然可用

## 相关文件

- `~/.zshrc` - Zsh 配置文件
- `~/.zprofile` - Zsh 登录配置文件
- `/opt/homebrew/opt/influxdb@1/bin/` - InfluxDB 可执行文件目录
