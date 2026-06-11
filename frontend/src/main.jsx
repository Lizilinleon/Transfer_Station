import React, { startTransition, useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  BookOpen,
  ChevronDown,
  CircleDollarSign,
  Copy,
  Database,
  Eye,
  Headphones,
  KeyRound,
  Megaphone,
  MessageSquare,
  Pencil,
  Plus,
  RefreshCcw,
  Search,
  Settings,
  Share2,
  Sparkles,
  Trash2,
  UserRound,
  Wallet,
  X
} from "lucide-react";
import "./styles/app.css";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
const DEFAULT_GROUP = import.meta.env.VITE_DEFAULT_GROUP || "codex-group";

const navItems = [
  { label: "密钥管理", path: "/app/keys" },
  { label: "体验中心", path: "/app/experience/chat" },
  { label: "模型服务", path: "/app/models" },
  { label: "文档", path: "/app/docs" },
  { label: "用户中心", path: "/app/user" }
];

const reservedFeatures = [
  { title: "对话体验", path: "/app/experience/chat", icon: MessageSquare, note: "预留聊天窗口和会话记录入口" },
  { title: "模型服务", path: "/app/models", icon: Database, note: "预留模型管理、价格、上下文配置" },
  { title: "渠道管理", path: "/app/channels", icon: Settings, note: "预留 DeepSeek / OpenAI 等渠道配置" },
  { title: "能力映射", path: "/app/abilities", icon: Sparkles, note: "预留模型到渠道的路由规则" },
  { title: "用量日志", path: "/app/usage-logs", icon: CircleDollarSign, note: "预留请求日志、额度和消耗统计" }
];

function apiUrl(path) {
  return `${API_BASE}${path}`;
}

async function requestJSON(path, options = {}) {
  const response = await fetch(apiUrl(path), {
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
    ...options
  });
  const result = await response.json().catch(() => null);

  if (!response.ok) {
    throw new Error(result?.message || result?.error?.message || `request failed: ${response.status}`);
  }

  return result;
}

function formatDate(value) {
  if (!value) {
    return "-";
  }
  return new Date(value).toLocaleString("zh-CN", { hour12: false });
}

function isPathActive(currentPath, itemPath) {
  if (itemPath === "/app/keys") {
    return currentPath === "/" || currentPath === "/app" || currentPath.startsWith("/app/keys");
  }
  return currentPath.startsWith(itemPath);
}

function TopNav({ currentPath }) {
  return (
    <header className="top-nav">
      <div className="brand">
        <span className="brand-mark" />
        <span>控制台</span>
      </div>
      <nav className="nav-links">
        {navItems.map((item) => (
          <a className={isPathActive(currentPath, item.path) ? "active" : ""} href={item.path} key={item.label}>
            {item.label}
          </a>
        ))}
      </nav>
      <div className="nav-actions">
        <button className="ghost-btn" type="button">
          <BookOpen size={18} /> 新手教程
        </button>
        <button className="ghost-btn" type="button">
          <Share2 size={18} /> 推广
        </button>
        <a className="recharge-btn" href="/app/recharge">
          <Wallet size={18} /> 充值
        </a>
        <button className="icon-btn" aria-label="消息" type="button">
          <MessageSquare size={18} />
        </button>
        <button className="icon-btn" aria-label="设置" type="button">
          <Settings size={18} />
        </button>
        <div className="user-chip">
          <span>测试用户</span>
          <ChevronDown size={16} />
        </div>
      </div>
    </header>
  );
}

function AnnouncementBar() {
  return (
    <div className="announcement">
      <span className="announce-icon">
        <Megaphone size={16} />
      </span>
      <strong>公告</strong>
      <span className="announce-text">
        当前先完成前后端联通：密钥管理读取 `/api/models` 和 `/api/keys`，其他模块先放置入口。
      </span>
    </div>
  );
}

function SideBar({ currentPath }) {
  return (
    <aside className="page-sidebar">
      <div className="sidebar-company">
        <div className="sidebar-company-icon">
          <Database size={20} />
        </div>
        <div>
          <h2>示例团队</h2>
          <p>开发环境</p>
        </div>
      </div>

      <div className="sidebar-section-title">当前功能</div>
      <a className={isPathActive(currentPath, "/app/keys") ? "sidebar-menu active" : "sidebar-menu"} href="/app/keys">
        <span className="sidebar-menu-icon">
          <KeyRound size={19} />
        </span>
        <span>密钥管理</span>
        <em className="sidebar-status">已接入</em>
      </a>

      <div className="sidebar-section-title">功能预留</div>
      {reservedFeatures.slice(1).map((item) => {
        const Icon = item.icon;
        return (
          <a
            className={isPathActive(currentPath, item.path) ? "sidebar-menu sidebar-menu-muted active" : "sidebar-menu sidebar-menu-muted"}
            href={item.path}
            key={item.path}
          >
            <span className="sidebar-menu-icon sidebar-menu-icon-muted">
              <Icon size={18} />
            </span>
            <span>{item.title}</span>
            <em className="sidebar-status">预留</em>
          </a>
        );
      })}

      <button className="sidebar-collapse" type="button">
        <ChevronDown size={18} />
        <span>收起</span>
      </button>
    </aside>
  );
}

function BackendStatus({ status, info, errorMessage }) {
  const labelMap = {
    checking: "连接中",
    online: "已连接",
    offline: "未连接"
  };

  return (
    <div className="backend-strip">
      <div>
        <span className={`backend-chip ${status}`}>{labelMap[status]}</span>
        <strong>后端服务</strong>
      </div>
      <p>
        {status === "online"
          ? `${info?.app_name || "AI Smart Chat Platform"} · ${info?.env || "development"} · ${API_BASE}`
          : errorMessage || `请先启动 Go 服务：${API_BASE}`}
      </p>
    </div>
  );
}

function ModelSelector({ availableModels, selectedModels, onToggleModel }) {
  return (
    <div className="model-selector">
      {availableModels.map((modelName) => {
        const active = selectedModels.includes(modelName);
        return (
          <button
            className={active ? "model-chip active" : "model-chip"}
            key={modelName}
            onClick={() => onToggleModel(modelName)}
            type="button"
          >
            {modelName}
          </button>
        );
      })}
    </div>
  );
}

function KeyEditor({ draft, title, availableModels, onChange, onToggleModel, onSubmit, onCancel }) {
  return (
    <section className="editor-card">
      <div className="editor-head">
        <div>
          <h3>{title}</h3>
          <p>保存后会直接写入后端 SQLite，后续 `/v1/chat/completions` 可使用这里生成的密钥。</p>
        </div>
        <button className="table-icon-btn" onClick={onCancel} type="button" aria-label="关闭编辑面板">
          <X size={18} />
        </button>
      </div>

      <div className="editor-grid">
        <label>
          <span>密钥名称</span>
          <input value={draft.name} onChange={(event) => onChange("name", event.target.value)} />
        </label>
        <label>
          <span>总额度</span>
          <input
            type="number"
            min="0"
            value={draft.quota}
            onChange={(event) => onChange("quota", Number(event.target.value))}
          />
        </label>
        <label className="editor-wide">
          <span>可用模型</span>
          <ModelSelector
            availableModels={availableModels}
            selectedModels={draft.modelNames}
            onToggleModel={onToggleModel}
          />
        </label>
      </div>

      <div className="editor-actions">
        <button className="outline-action" onClick={onCancel} type="button">
          取消
        </button>
        <button className="create-solid-btn" onClick={onSubmit} type="button">
          <Plus size={18} />
          保存密钥
        </button>
      </div>
    </section>
  );
}

function KeyTable({ rows, onToggleEnabled, onDelete, onEdit, onOpenModelPicker }) {
  if (!rows.length) {
    return (
      <section className="token-table-card empty-state">
        <div className="empty-copy">
          <h3>暂无密钥</h3>
          <p>后端数据库里还没有密钥，可以点击“创建密钥”写入第一条数据。</p>
        </div>
      </section>
    );
  }

  return (
    <section className="token-table-card">
      <div className="token-table-head">
        <span>名称</span>
        <span>密钥</span>
        <span>可用模型</span>
        <span>状态</span>
        <span>额度</span>
        <span>创建时间</span>
        <span>操作</span>
      </div>

      {rows.map((row) => (
        <div className="token-table-row" key={row.id}>
          <span className="token-name">{row.name}</span>
          <div className="token-secret-wrap">
            <span className="token-secret">{row.value}</span>
            <button className="table-icon-btn" aria-label="查看密钥" type="button">
              <Eye size={16} />
            </button>
            <button
              className="table-icon-btn"
              aria-label="复制密钥"
              onClick={() => navigator.clipboard?.writeText(row.value)}
              type="button"
            >
              <Copy size={16} />
            </button>
          </div>
          <button className="table-pill table-pill-button" onClick={() => onOpenModelPicker(row.id)} type="button">
            <Database size={16} />
            {row.modelNames.length ? `${row.modelNames.length} 个模型` : "选择模型"}
          </button>
          <span className={row.enabled ? "status-pill" : "status-pill status-pill-disabled"}>
            {row.enabled ? "已启用" : "已禁用"}
          </span>
          <span className="token-quota">
            已用 {row.usedQuota.toFixed(2)}
            <small>剩余 {row.remainingQuota.toFixed(2)} / 总额 {row.quota}</small>
          </span>
          <span className="token-date">{row.createdAt}</span>
          <div className="row-actions">
            <button className="action-outline" aria-label="编辑" onClick={() => onEdit(row.id)} type="button">
              <Pencil size={16} />
            </button>
            <button className="action-warning" onClick={() => onToggleEnabled(row)} type="button">
              <Eye size={15} />
              {row.enabled ? "禁用" : "启用"}
            </button>
            <button className="action-danger" aria-label="删除" onClick={() => onDelete(row.id)} type="button">
              <Trash2 size={16} />
            </button>
          </div>
        </div>
      ))}

      <div className="token-table-footer">
        <span>
          显示第 1 条 - 第 {rows.length} 条，共 {rows.length} 条
        </span>
        <div className="pagination">
          <button className="page-arrow" aria-label="上一页" type="button">
            ‹
          </button>
          <span className="page-number active">1</span>
          <button className="page-arrow" aria-label="下一页" type="button">
            ›
          </button>
        </div>
      </div>
    </section>
  );
}

function ReservedFeatureGrid() {
  return (
    <section className="reserved-grid">
      {reservedFeatures.map((feature) => {
        const Icon = feature.icon;
        return (
          <a className="reserved-card" href={feature.path} key={feature.path}>
            <Icon size={22} />
            <strong>{feature.title}</strong>
            <span>{feature.note}</span>
          </a>
        );
      })}
    </section>
  );
}

function mapKeyItem(item) {
  return {
    id: item.id,
    name: item.name,
    value: item.access_key,
    groupName: item.group_name || DEFAULT_GROUP,
    modelNames: item.model_names || [],
    enabled: item.enabled,
    quota: item.quota || 0,
    usedQuota: item.used_quota || 0,
    remainingQuota: item.remaining_quota || 0,
    unlimitedQuota: item.unlimited_quota || false,
    createdAt: formatDate(item.created_at),
    billingRule: item.billing_rule || "reserved",
    remark: item.remark || ""
  };
}

function KeysPage({ currentPath }) {
  const [rows, setRows] = useState([]);
  const [availableModels, setAvailableModels] = useState(["deepseek-chat"]);
  const [searchText, setSearchText] = useState("");
  const [showEditor, setShowEditor] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(false);
  const [backendStatus, setBackendStatus] = useState("checking");
  const [backendInfo, setBackendInfo] = useState(null);
  const [errorMessage, setErrorMessage] = useState("");
  const [draft, setDraft] = useState({
    name: "",
    quota: 100,
    modelNames: ["deepseek-chat"]
  });

  const selectedModelFallback = useMemo(() => [availableModels[0] || "deepseek-chat"], [availableModels]);

  async function fetchHealth() {
    const result = await requestJSON("/health");
    setBackendInfo(result?.data || null);
    setBackendStatus("online");
  }

  async function fetchModels() {
    const result = await requestJSON("/api/models");
    const items = result?.data?.items || [];
    if (!items.length) {
      return ["deepseek-chat"];
    }
    return items.map((item) => item.name);
  }

  async function fetchKeys(keyword = "") {
    const query = keyword ? `?keyword=${encodeURIComponent(keyword)}` : "";
    const result = await requestJSON(`/api/keys${query}`);
    return (result?.data?.items || []).map(mapKeyItem);
  }

  async function reloadData(keyword = searchText) {
    setLoading(true);
    setErrorMessage("");

    try {
      await fetchHealth();
      const [models, keys] = await Promise.all([fetchModels(), fetchKeys(keyword)]);
      setAvailableModels(models);
      setRows(keys);
    } catch (error) {
      setBackendStatus("offline");
      setErrorMessage("后端连接失败，请确认 Go 服务已经启动。");
      console.error(error);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    reloadData("");
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      reloadData(searchText);
    }, 250);

    return () => window.clearTimeout(timer);
  }, [searchText]);

  function resetDraft(name = "", models = selectedModelFallback) {
    setDraft({
      name,
      quota: 100,
      modelNames: models.length ? models : ["deepseek-chat"]
    });
  }

  function handleDraftChange(field, value) {
    setDraft((currentDraft) => ({
      ...currentDraft,
      [field]: value
    }));
  }

  function handleToggleDraftModel(modelName) {
    setDraft((currentDraft) => {
      const exists = currentDraft.modelNames.includes(modelName);
      return {
        ...currentDraft,
        modelNames: exists
          ? currentDraft.modelNames.filter((item) => item !== modelName)
          : [...currentDraft.modelNames, modelName]
      };
    });
  }

  function handleOpenCreator() {
    setEditingId(null);
    resetDraft(`测试密钥 ${rows.length + 1}`, selectedModelFallback);
    setShowEditor(true);
  }

  function handleEdit(id) {
    const current = rows.find((row) => row.id === id);
    if (!current) {
      return;
    }

    setEditingId(id);
    setDraft({
      name: current.name,
      quota: current.quota,
      modelNames: current.modelNames.length ? current.modelNames : selectedModelFallback
    });
    setShowEditor(true);
  }

  function handleOpenModelPicker(id) {
    handleEdit(id);
  }

  function buildKeyPayload(currentRow = null, enabled = true) {
    return {
      name: draft.name.trim(),
      group_name: currentRow?.groupName || DEFAULT_GROUP,
      model_names: draft.modelNames,
      enabled,
      quota: draft.quota,
      used_quota: currentRow?.usedQuota || 0,
      unlimited_quota: currentRow?.unlimitedQuota || false,
      billing_rule: currentRow?.billingRule || "reserved",
      billing_config: "",
      input_token_price: 0,
      output_token_price: 0,
      request_price: 0,
      remark: currentRow?.remark || "frontend integration write"
    };
  }

  async function handleSaveDraft() {
    if (!draft.name.trim()) {
      setErrorMessage("请先填写密钥名称。");
      return;
    }

    try {
      if (editingId !== null) {
        const current = rows.find((row) => row.id === editingId);
        await requestJSON(`/api/keys/${editingId}`, {
          method: "PUT",
          body: JSON.stringify(buildKeyPayload(current, current ? current.enabled : true))
        });
      } else {
        await requestJSON("/api/keys", {
          method: "POST",
          body: JSON.stringify(buildKeyPayload(null, true))
        });
      }

      startTransition(() => {
        setEditingId(null);
        setShowEditor(false);
        resetDraft();
      });
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("保存密钥失败，请检查后端接口是否可用。");
      console.error(error);
    }
  }

  function handleCancelEditor() {
    setEditingId(null);
    setShowEditor(false);
    resetDraft();
  }

  async function handleToggleEnabled(row) {
    try {
      await requestJSON(`/api/keys/${row.id}`, {
        method: "PUT",
        body: JSON.stringify({
          name: row.name,
          group_name: row.groupName || DEFAULT_GROUP,
          model_names: row.modelNames,
          enabled: !row.enabled,
          quota: row.quota,
          used_quota: row.usedQuota,
          unlimited_quota: row.unlimitedQuota,
          billing_rule: row.billingRule || "reserved",
          billing_config: "",
          input_token_price: 0,
          output_token_price: 0,
          request_price: 0,
          remark: row.remark || "status toggle"
        })
      });
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("切换状态失败。");
      console.error(error);
    }
  }

  async function handleDelete(id) {
    try {
      await requestJSON(`/api/keys/${id}`, { method: "DELETE" });
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("删除密钥失败。");
      console.error(error);
    }
  }

  return (
    <main className="keys-layout">
      <SideBar currentPath={currentPath} />

      <section className="keys-main">
        <div className="page-title-row">
          <div className="title-left">
            <h1>密钥管理</h1>
            <p>先接通后端数据，后续再逐步补齐模型、渠道、日志等模块。</p>
          </div>
          <div className="title-actions">
            <button className="outline-action" onClick={() => reloadData(searchText)} type="button">
              <RefreshCcw size={18} />
              刷新
            </button>
            <button className="create-solid-btn" onClick={handleOpenCreator} type="button">
              <Plus size={18} />
              创建密钥
            </button>
          </div>
        </div>

        <BackendStatus status={backendStatus} info={backendInfo} errorMessage={errorMessage} />

        <div className="info-banner">
          <KeyRound size={16} />
          <span>
            当前页面已接入 `{API_BASE}/api/models`、`{API_BASE}/api/keys`，默认写入分组 `{DEFAULT_GROUP}`。
          </span>
        </div>

        {showEditor ? (
          <KeyEditor
            draft={draft}
            title={editingId !== null ? "编辑密钥" : "创建密钥"}
            availableModels={availableModels}
            onChange={handleDraftChange}
            onToggleModel={handleToggleDraftModel}
            onSubmit={handleSaveDraft}
            onCancel={handleCancelEditor}
          />
        ) : null}

        <div className="search-inline">
          <Search size={19} />
          <input
            placeholder="搜索密钥名称、Key 或分组..."
            value={searchText}
            onChange={(event) => setSearchText(event.target.value)}
          />
        </div>

        {errorMessage ? <p className="inline-message error">{errorMessage}</p> : null}
        {loading ? <p className="inline-message">正在同步后端数据...</p> : null}

        <KeyTable
          rows={rows}
          onToggleEnabled={handleToggleEnabled}
          onDelete={handleDelete}
          onEdit={handleEdit}
          onOpenModelPicker={handleOpenModelPicker}
        />

        <div className="section-heading">
          <h2>后续功能入口</h2>
          <p>先把位置留出来，后面按模块逐个接后端接口。</p>
        </div>
        <ReservedFeatureGrid />
      </section>

      <button className="float-help" aria-label="联系客服" type="button">
        <Headphones size={34} />
        <span />
      </button>
    </main>
  );
}

function PlaceholderPage({ title, icon: Icon, description }) {
  return (
    <main className="placeholder-page">
      <div className="placeholder-card">
        <Icon size={42} />
        <h1>{title}</h1>
        <p>{description}</p>
        <div className="placeholder-actions">
          <a className="create-solid-btn" href="/app/keys">
            返回密钥管理
          </a>
          <button className="outline-action" type="button">
            暂不实现
          </button>
        </div>
      </div>
    </main>
  );
}

function App() {
  const path = window.location.pathname;
  let page = <KeysPage currentPath={path} />;

  if (path.includes("/experience/chat")) {
    page = (
      <PlaceholderPage
        title="体验中心 · 对话"
        icon={MessageSquare}
        description="这里预留聊天窗口、模型选择、会话列表和消息记录，下一步可以接 `/v1/chat/completions`。"
      />
    );
  }
  if (path.includes("/experience/image")) {
    page = (
      <PlaceholderPage
        title="图片生成"
        icon={Sparkles}
        description="这里预留图片生成表单、任务列表和结果展示区。"
      />
    );
  }
  if (path.includes("/models")) {
    page = (
      <PlaceholderPage
        title="模型服务"
        icon={Database}
        description="这里预留模型 CRUD、价格配置、上下文长度和默认模型开关。"
      />
    );
  }
  if (path.includes("/channels")) {
    page = (
      <PlaceholderPage
        title="渠道管理"
        icon={Settings}
        description="这里预留上游渠道、Base URL、API Key、权重和优先级配置。"
      />
    );
  }
  if (path.includes("/abilities")) {
    page = (
      <PlaceholderPage
        title="能力映射"
        icon={Sparkles}
        description="这里预留 group + model 到 channel 的路由规则管理。"
      />
    );
  }
  if (path.includes("/usage-logs")) {
    page = (
      <PlaceholderPage
        title="用量日志"
        icon={CircleDollarSign}
        description="这里预留请求记录、token 消耗、额度扣减和错误统计。"
      />
    );
  }
  if (path.includes("/docs")) {
    page = (
      <PlaceholderPage
        title="文档中心"
        icon={BookOpen}
        description="这里预留接口文档、快速开始和示例请求。"
      />
    );
  }
  if (path.includes("/user")) {
    page = (
      <PlaceholderPage
        title="用户中心"
        icon={UserRound}
        description="这里预留用户资料、团队信息和账号设置。"
      />
    );
  }
  if (path.includes("/recharge")) {
    page = (
      <PlaceholderPage
        title="充值中心"
        icon={CircleDollarSign}
        description="这里预留余额、套餐、订单和支付回调状态。"
      />
    );
  }

  return (
    <div className="app">
      <TopNav currentPath={path} />
      <AnnouncementBar />
      {page}
    </div>
  );
}

createRoot(document.getElementById("root")).render(<App />);
