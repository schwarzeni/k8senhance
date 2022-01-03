未实现：

- 实现 manifest 缓存，以及其它节点获取 manifest 缓存的功能 [finish]
- 代理下载 layer 请求时策略的实现（参考飞书流程图，目前仅仅是随机选取一个存在 layer 的节点）
- 使用 cache.Layer 来验证 layer 是否存在，同时在代理的时候更新
- imagecachecontroller 实现功能：定时移除未及时上传心跳的节点信息
- 重构一下节点数据是由那个结构体记录的，from cache.metric or response or cache.nodeinfos
- 研究 response 需要返回那些标准性的 header （目前阶段还需要请求一下远端的仓库）

---

仅支持 sha256 格式的 layer
