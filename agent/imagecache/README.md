未实现：

- 实现 manifest 缓存，以及其它节点获取 manifest 缓存的功能
- 代理下载 layer 请求时策略的实现（参考飞书流程图，目前仅仅是随机选取一个存在 layer 的节点）
- 使用 cache.Layer 来验证 layer 是否存在，同时在代理的时候更新
- imagecachecontroller 实现功能：定时移除未及时上传心跳的节点信息
- 重构一下节点数据是由那个结构体记录的，from cache.metric or response or cache.nodeinfos
