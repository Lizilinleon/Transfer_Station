import React, { startTransition, useEffect, useState } from "react";
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

const API_BASE = "http://localhost:8080";

const navItems = [
  { label: "秘钥管理", path: "/app/keys", active: true },
  { label: "体验中心", path: "/app/experience/chat" },
  { label: "模型服务", path: "/app/models" },
  { label: "文档", path: "/app/docs" },
  { label: "用户中心", path: "/app/user" }
];

function TopNav() {
  return (
    <header className="top-nav">
      <div className="brand">
        <span className="brand-mark" />
        <span>控制台</span>
      </div>
      <nav className="nav-links">
        {navItems.map((item) => (
          <a className={item.active ? "active" : ""} href={item.path} key={item.label}>
            {item.label}
          </a>
        ))}
      </nav>
      <div className="nav-actions">
        <button className="ghost-btn"><BookOpen size={18} /> 新手教程</button>
        <button className="ghost-btn"><Share2 size={18} /> 推广</button>
        <a className="recharge-btn" href="/app/recharge"><Wallet size={18} /> 充值</a>
        <button className="icon-btn" aria-label="消息"><MessageSquare size={18} /></button>
        <button className="icon-btn" aria-label="设置"><Settings size={18} /></button>
        <div className="user-chip"><span>测试用户</span><ChevronDown size={16} /></div>
      </div>
    </header>
  );
}

function AnnouncementBar() {
  return (
    <div className="announcement">
      <span className="announce-icon"><Megaphone size={16} /></span>
      <strong>公告</strong>
      <span className="announce-text">本平台提供稳定的 AI 服务测试环境，企业合作请联系客服。</span>
    </div>
  );
}

function SideBar() {
  return (
    <aside className="page-sidebar">
      <div className="sidebar-company">
        <div className="sidebar-company-icon">
          <Database size={20} />
        </div>
        <div>
          <h2>示例团队</h2>
          <p>成员</p>
        </div>
      </div>

      <button className="sidebar-menu active">
        <span className="sidebar-menu-icon">
          <KeyRound size={19} />
        </span>
        <span>令牌管理</span>
      </button>

      <button className="sidebar-collapse">
        <ChevronDown size={18} />
        <span>收起</span>
      </button>
    </aside>
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
          <p>这里已经接入后端字段结构，保存后会直接写入 SQLite 数据库。</p>
        </div>
        <button className="table-icon-btn" onClick={onCancel} type="button" aria-label="关闭编辑面板">
          <X size={18} />
        </button>
      </div>

      <div className="editor-grid">
        <label>
          <span>令牌名称</span>
          <input value={draft.name} onChange={(event) => onChange("name", event.target.value)} />
        </label>
        <label>
          <span>额度</span>
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
        <button className="outline-action" onClick={onCancel} type="button">取消</button>
        <button className="create-solid-btn" onClick={onSubmit} type="button">
          <Plus size={18} />
          保存令牌
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
          <h3>暂无令牌</h3>
          <p>数据库里还没有数据，可以直接点击创建令牌写入后端。</p>
        </div>
      </section>
    );
  }

  return (
    <section className="token-table-card">
      <div className="token-table-head">
        <span>名称</span>
        <span>令牌密钥</span>
        <span>可用模型</span>
        <span>状态</span>
        <span>已用额度</span>
        <span>创建时间</span>
        <span>操作</span>
      </div>

      {rows.map((row) => (
        <div className="token-table-row" key={row.id}>
          <span className="token-name">{row.name}</span>
          <div className="token-secret-wrap">
            <span className="token-secret">{row.value}</span>
            <button className="table-icon-btn" aria-label="查看令牌"><Eye size={16} /></button>
            <button
              className="table-icon-btn"
              aria-label="复制令牌"
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
            ¥{row.usedQuota.toFixed(2)}
            <small>/ 总额度 {row.quota}</small>
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
        <span>显示第 1 条-第 {rows.length} 条，共 {rows.length} 条</span>
        <div className="pagination">
          <button className="page-arrow" aria-label="上一页" type="button">‹</button>
          <span className="page-number active">1</span>
          <button className="page-arrow" aria-label="下一页" type="button">›</button>
        </div>
      </div>
    </section>
  );
}

function mapKeyItem(item) {
  return {
    id: item.id,
    name: item.name,
    value: item.access_key,
    modelNames: item.model_names || [],
    enabled: item.enabled,
    quota: item.quota || 0,
    usedQuota: item.used_quota || 0,
    createdAt: item.created_at ? new Date(item.created_at).toLocaleString("zh-CN", { hour12: false }) : "-",
    billingRule: item.billing_rule || "reserved"
  };
}

function KeysPage() {
  const [rows, setRows] = useState([]);
  const [availableModels, setAvailableModels] = useState(["gpt-5.5"]);
  const [searchText, setSearchText] = useState("");
  const [showEditor, setShowEditor] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [draft, setDraft] = useState({
    name: "",
    quota: 100,
    modelNames: ["gpt-5.5"]
  });

  async function fetchModels() {
    const response = await fetch(`${API_BASE}/api/models`);
    const result = await response.json();
    const items = result?.data?.items || [];
    if (!items.length) {
      return ["gpt-5.5"];
    }
    return items.map((item) => item.name);
  }

  async function fetchKeys(keyword = "") {
    const query = keyword ? `?keyword=${encodeURIComponent(keyword)}` : "";
    const response = await fetch(`${API_BASE}/api/keys${query}`);
    const result = await response.json();
    return (result?.data?.items || []).map(mapKeyItem);
  }

  async function reloadData(keyword = searchText) {
    setLoading(true);
    setErrorMessage("");

    try {
      const [models, keys] = await Promise.all([fetchModels(), fetchKeys(keyword)]);
      setAvailableModels(models);
      setRows(keys);
    } catch (error) {
      setErrorMessage("后端连接失败，请确认 Go 服务已启动。");
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

  function resetDraft(name = "", models = ["gpt-5.5"]) {
    setDraft({
      name,
      quota: 100,
      modelNames: models.length ? models : ["gpt-5.5"]
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
    resetDraft(`测试令牌 ${rows.length + 1}`, [availableModels[0] || "gpt-5.5"]);
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
      modelNames: current.modelNames.length ? current.modelNames : [availableModels[0] || "gpt-5.5"]
    });
    setShowEditor(true);
  }

  function handleOpenModelPicker(id) {
    handleEdit(id);
  }

  async function handleSaveDraft() {
    if (!draft.name.trim()) {
      return;
    }

    const payload = {
      name: draft.name.trim(),
      group_name: "codex专用分组",
      model_names: draft.modelNames,
      enabled: true,
      quota: draft.quota,
      used_quota: 0,
      unlimited_quota: false,
      billing_rule: "reserved",
      billing_config: "",
      input_token_price: 0,
      output_token_price: 0,
      request_price: 0,
      remark: "前端联调写入"
    };

    try {
      if (editingId !== null) {
        const current = rows.find((row) => row.id === editingId);
        const response = await fetch(`${API_BASE}/api/keys/${editingId}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            ...payload,
            enabled: current ? current.enabled : true,
            used_quota: current ? current.usedQuota : 0
          })
        });
        if (!response.ok) {
          throw new Error("update key failed");
        }
      } else {
        const response = await fetch(`${API_BASE}/api/keys`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload)
        });
        if (!response.ok) {
          throw new Error("create key failed");
        }
      }

      startTransition(() => {
        setEditingId(null);
        setShowEditor(false);
        resetDraft();
      });
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("保存令牌失败，请检查后端是否可用。");
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
      const response = await fetch(`${API_BASE}/api/keys/${row.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: row.name,
          group_name: "codex专用分组",
          model_names: row.modelNames,
          enabled: !row.enabled,
          quota: row.quota,
          used_quota: row.usedQuota,
          unlimited_quota: false,
          billing_rule: row.billingRule || "reserved",
          billing_config: "",
          input_token_price: 0,
          output_token_price: 0,
          request_price: 0,
          remark: "状态切换"
        })
      });
      if (!response.ok) {
        throw new Error("toggle key failed");
      }
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("切换状态失败。");
      console.error(error);
    }
  }

  async function handleDelete(id) {
    try {
      const response = await fetch(`${API_BASE}/api/keys/${id}`, {
        method: "DELETE"
      });
      if (!response.ok) {
        throw new Error("delete key failed");
      }
      reloadData(searchText);
    } catch (error) {
      setErrorMessage("删除令牌失败。");
      console.error(error);
    }
  }

  return (
    <main className="keys-layout">
      <SideBar />

      <section className="keys-main">
        <div className="page-title-row">
          <div className="title-left">
            <h1>令牌管理</h1>
            <p>管理您创建的 API 令牌</p>
          </div>
          <div className="title-actions">
            <button className="outline-action" onClick={() => reloadData(searchText)} type="button">
              <RefreshCcw size={18} />
              刷新
            </button>
            <button className="create-solid-btn" onClick={handleOpenCreator} type="button">
              <Plus size={18} />
              创建令牌
            </button>
          </div>
        </div>

        <div className="info-banner">
          <KeyRound size={16} />
          <span>当前页面已接入后端 `/api/models` 和 `/api/keys`，创建、搜索、编辑、删除都会直接落到 SQLite。</span>
        </div>

        {showEditor ? (
          <KeyEditor
            draft={draft}
            title={editingId !== null ? "编辑令牌" : "创建令牌"}
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
            placeholder="搜索令牌名称、Key 或模型..."
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
      </section>

      <button className="float-help" aria-label="联系客服"><Headphones size={34} /><span /></button>
    </main>
  );
}

function PlaceholderPage({ title, icon: Icon }) {
  return (
    <main className="placeholder-page">
      <div className="placeholder-card">
        <Icon size={42} />
        <h1>{title}</h1>
        <p>这里先保留静态占位，后续再逐步接入具体功能。</p>
      </div>
    </main>
  );
}

function App() {
  const path = window.location.pathname;
  let page = <KeysPage />;

  if (path.includes("/experience/chat")) {
    page = <PlaceholderPage title="体验中心 · 对话" icon={MessageSquare} />;
  }
  if (path.includes("/experience/image")) {
    page = <PlaceholderPage title="图片生成" icon={Sparkles} />;
  }
  if (path.includes("/models")) {
    page = <PlaceholderPage title="模型服务" icon={Database} />;
  }
  if (path.includes("/docs")) {
    page = <PlaceholderPage title="文档中心" icon={BookOpen} />;
  }
  if (path.includes("/user")) {
    page = <PlaceholderPage title="用户中心" icon={UserRound} />;
  }
  if (path.includes("/recharge")) {
    page = <PlaceholderPage title="充值中心" icon={CircleDollarSign} />;
  }

  return (
    <div className="app">
      <TopNav />
      <AnnouncementBar />
      {page}
    </div>
  );
}

createRoot(document.getElementById("root")).render(<App />);
