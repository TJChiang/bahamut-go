# 巴哈姆特自動化簽到

Github 和 Gitlab 有很多開發者做過這個專案，自己希望以 Go 實現這個服務

## Todo

- [x] 登入模組
- [x] 簽到模組
- [x] 設定 GitHub Actions
- [x] Docker化
- [x] 設定 CronJob
- [x] 串接 Notify
- [ ] 使用文件
- [ ] 寫測試
- [ ] Docker image 優化
- [ ] GitHub Actions 優化

## Troubleshooting

- Docker 出現 npx install playwright 錯誤：[run playwright-go on Ubuntu](https://github.com/playwright-community/playwright-go/issues/277)
- nektos/act 本機測試 Github Action 時出現 `Cannot connect to Docker daemon` (colima) => [解決辦法](https://github.com/nektos/act/issues/1051)
- Github Action Cache Docker image => [Solution](https://stackoverflow.com/a/71183339)