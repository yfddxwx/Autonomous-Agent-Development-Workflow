# Code Quality & Testing Guide
name: code-quality
description: 保证代码质量、模块解耦与单元测试覆盖

## 核心要求:
1. 每次编写或修改 Go 代码时，必须保持函数职责单一，避免超长函数。
2. 编码完成后必须执行 go test 验证正确性。
3. 遵循 Clean Architecture 架构分层规范，禁止跨层循环依赖。
