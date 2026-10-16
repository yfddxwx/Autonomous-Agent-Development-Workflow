// DevAgent Web UI Client Logic (Google Material 3 / Gemini Theme)
document.addEventListener('DOMContentLoaded', () => {
  let ws = null;
  let pendingApprovalToken = null;
  let activeProfileId = 'deepseek';
  let activeSessionId = '';

  // DOM Elements - Main Layout
  const messagesContainer = document.getElementById('messages-container');
  const welcomeHero = document.getElementById('welcome-hero');
  const promptInput = document.getElementById('prompt-input');
  const btnSend = document.getElementById('btn-send');
  const btnStop = document.getElementById('btn-stop');
  const statusDot = document.getElementById('status-dot');
  const statusText = document.getElementById('status-text');
  const workspacePath = document.getElementById('workspace-path');
  const inputWsHint = document.getElementById('input-ws-hint');
  const btnSwitchWorkspace = document.getElementById('btn-switch-workspace');
  const modelSelect = document.getElementById('model-select');
  const approvalModeSelect = document.getElementById('approval-mode-select');
  const sessionSnapshotSelect = document.getElementById('session-snapshot-select');
  const btnResetSession = document.getElementById('btn-reset-session');
  const btnPromptResetSession = document.getElementById('btn-prompt-reset-session');
  const sessionTokenCount = document.getElementById('session-token-count');

  // Thinking Level Elements (Pi Agent Style Reasoning Depth)
  const btnThinkingLevel = document.getElementById('btn-thinking-level');
  const thinkingLevelLabel = document.getElementById('thinking-level-label');
  const thinkingLevelIcon = document.getElementById('thinking-level-icon');
  const thinkingLevelChevron = document.getElementById('thinking-level-chevron');
  const thinkingLevelMenu = document.getElementById('thinking-level-menu');
  const thinkingLevelOptions = document.getElementById('thinking-level-options');
  const thinkingModelName = document.getElementById('thinking-model-name');
  let currentThinkingLevel = localStorage.getItem('devagent_thinking_level') || 'high';
  let allThinkingDefinitions = {
    off: { level: 'off', label: '关闭思考', tokens: '0', desc: '不进行推理 (No reasoning)' },
    minimal: { level: 'minimal', label: '极简思考', tokens: '~1k', desc: '极短推理 (~1k tokens)' },
    low: { level: 'low', label: '浅度思考', tokens: '~2k', desc: '轻度推理 (~2k tokens)' },
    medium: { level: 'medium', label: '中度思考', tokens: '~8k', desc: '中度推理 (~8k tokens)' },
    high: { level: 'high', label: '深度思考', tokens: '~16k', desc: '深度推理 (~16k tokens)' },
    xhigh: { level: 'xhigh', label: '超深思考', tokens: '~32k', desc: '超高推理 (~32k tokens)' },
    max: { level: 'max', label: '极限思考', tokens: 'Max', desc: '最大算力推理 (Max reasoning)' }
  };
  let cachedProfilesData = [];

  // Git Widget Elements
  const gitWidget = document.getElementById('git-widget');
  const gitStatusText = document.getElementById('git-status-text');
  const btnGitRollback = document.getElementById('btn-git-rollback');

  // Telemetry Elements
  const telemetryTokens = document.getElementById('telemetry-tokens');
  const telemetryCost = document.getElementById('telemetry-cost');

  // Workspace Modal Elements
  const workspaceModal = document.getElementById('workspace-modal');
  const btnCloseWsModal = document.getElementById('btn-close-ws-modal');
  const btnCancelWs = document.getElementById('btn-cancel-ws');
  const btnConfirmSwitchWs = document.getElementById('btn-confirm-switch-ws');
  const wsNewPath = document.getElementById('ws-new-path');
  const wsSwitchError = document.getElementById('ws-switch-error');
  const btnBrowseFolder = document.getElementById('btn-browse-folder');
  const btnBrowseFolderText = document.getElementById('btn-browse-folder-text');

  // Approval Elements
  const approvalCard = document.getElementById('approval-card');
  const approvalToolBadge = document.getElementById('approval-tool-badge');
  const approvalDesc = document.getElementById('approval-desc');
  const approvalPayload = document.getElementById('approval-payload');
  const approvalRejectReason = document.getElementById('approval-reject-reason');
  const btnRejectApproval = document.getElementById('btn-reject-approval');
  const btnConfirmApproval = document.getElementById('btn-confirm-approval');

  // Settings Modal Elements
  const btnSettings = document.getElementById('btn-settings');
  const settingsModal = document.getElementById('settings-modal');
  const btnCloseSettings = document.getElementById('btn-close-settings');
  const btnDoneSettings = document.getElementById('btn-done-settings');
  const btnSaveSettings = document.getElementById('btn-save-settings');
  const activeProfileBadge = document.getElementById('active-profile-badge');
  const cfgBaseUrl = document.getElementById('cfg-base-url');
  const cfgApiKey = document.getElementById('cfg-api-key');
  const cfgModelName = document.getElementById('cfg-model-name');
  const cfgMaxSteps = document.getElementById('cfg-max-steps');

  // Model Profiles Manager Elements
  const modelProfilesContainer = document.getElementById('model-profiles-container');
  const btnToggleAddModel = document.getElementById('btn-toggle-add-model');
  const addModelPanel = document.getElementById('add-model-panel');
  const btnCancelAddModel = document.getElementById('btn-cancel-add-model');
  const btnConfirmAddModel = document.getElementById('btn-confirm-add-model');
  const newModelId = document.getElementById('new-model-id');
  const newModelName = document.getElementById('new-model-name');
  const newModelUrl = document.getElementById('new-model-url');
  const newModelKey = document.getElementById('new-model-key');
  const newModelCode = document.getElementById('new-model-code');
  const newModelSteps = document.getElementById('new-model-steps');

  // Session Drawer Elements
  const btnToggleDrawer = document.getElementById('btn-toggle-drawer');
  const sessionDrawer = document.getElementById('session-drawer');
  const btnCloseDrawer = document.getElementById('btn-close-drawer');
  const drawerBackdrop = document.getElementById('drawer-backdrop');
  const btnDrawerNewSession = document.getElementById('btn-drawer-new-session');
  const sessionList = document.getElementById('session-list');

  // Provider Hub Elements (Pi Agent Style)
  const providerCards = document.querySelectorAll('.btn-provider-card');
  const providerDiscoveryPanel = document.getElementById('provider-discovery-panel');
  const discProviderBadge = document.getElementById('disc-provider-badge');
  const btnCloseDiscPanel = document.getElementById('btn-close-disc-panel');
  const discBaseUrl = document.getElementById('disc-base-url');
  const discApiKey = document.getElementById('disc-api-key');
  const discStatusHint = document.getElementById('disc-status-hint');
  const btnFetchProviderModels = document.getElementById('btn-fetch-provider-models');
  const discModelsResult = document.getElementById('disc-models-result');
  const discFoundCount = document.getElementById('disc-found-count');
  const discSelectAll = document.getElementById('disc-select-all');
  const discModelsList = document.getElementById('disc-models-list');
  const btnImportDiscModels = document.getElementById('btn-import-disc-models');
  let selectedProvider = 'deepseek';
  let fetchedDiscoveredModels = [];

  function updateWorkspaceUI(ws) {
    if (!ws) return;
    const parts = ws.split(/[\\/]/).filter(Boolean);
    const base = parts.length > 0 ? parts[parts.length - 1] : ws;
    workspacePath.textContent = base;
    workspacePath.title = ws;
    if (btnSwitchWorkspace) {
      btnSwitchWorkspace.title = `当前工作空间: ${ws}\n点击更换目录`;
    }
    if (inputWsHint) {
      inputWsHint.textContent = base;
    }
  }

  // -------------------------------------------------------------
  // Drawer / Session Persistence Logic
  // -------------------------------------------------------------
  function openDrawer() {
    sessionDrawer.classList.remove('-translate-x-full');
    drawerBackdrop.classList.remove('hidden');
    loadSessions();
  }

  function closeDrawer() {
    sessionDrawer.classList.add('-translate-x-full');
    drawerBackdrop.classList.add('hidden');
  }

  if (btnToggleDrawer) btnToggleDrawer.addEventListener('click', openDrawer);
  if (btnCloseDrawer) btnCloseDrawer.addEventListener('click', closeDrawer);
  if (drawerBackdrop) drawerBackdrop.addEventListener('click', closeDrawer);

  let hasAutoLoadedActiveSession = false;

  async function loadSessions() {
    if (!sessionList) return;
    try {
      const res = await fetch('/api/sessions');
      if (!res.ok) return;
      const data = await res.json();
      const sessions = data.sessions || [];
      let targetId = data.active_id;
      if (!targetId && sessions.length > 0) {
        const withMessages = sessions.find(s => s.message_count > 0);
        targetId = withMessages ? withMessages.id : sessions[0].id;
      }

      if (targetId) {
        activeSessionId = targetId;
        if (!hasAutoLoadedActiveSession) {
          hasAutoLoadedActiveSession = true;
          switchToSession(targetId);
          return;
        }
      }

      sessionList.innerHTML = '';
      if (sessions.length === 0) {
        sessionList.innerHTML = '<div class="text-xs text-[#8e918f] text-center py-6">暂无历史会话记录</div>';
        return;
      }

      sessions.forEach(s => {
        const item = document.createElement('div');
        const isActive = s.id === activeSessionId;
        item.className = `p-3 rounded-2xl border transition cursor-pointer flex items-center justify-between group ${
          isActive
            ? 'bg-[#8ab4f8]/10 border-[#8ab4f8]/40 text-[#e3e3e3]'
            : 'bg-[#191a1b] hover:bg-[#282a2c] border-[#282a2c] text-[#c4c7c5]'
        }`;

        const formattedTime = new Date(s.updated_at).toLocaleString('zh-CN', {
          month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
        });

        const wsBase = s.workspace_dir ? (s.workspace_dir.split(/[\\/]/).filter(Boolean).pop() || s.workspace_dir) : '';
        item.innerHTML = `
          <div class="flex-1 min-w-0 pr-2">
            <div class="flex items-center gap-1.5">
              ${isActive ? '<span class="w-1.5 h-1.5 rounded-full bg-[#8ab4f8] flex-shrink-0"></span>' : ''}
              <div class="font-medium text-xs truncate ${isActive ? 'text-[#8ab4f8]' : 'text-[#e3e3e3]'}">${escapeHtml(s.title || '新研发会话')}</div>
            </div>
            <div class="flex items-center gap-2 mt-1 text-[10px] text-[#8e918f]">
              ${wsBase ? `<span class="text-[#8ab4f8] font-mono flex items-center gap-1 max-w-[110px] truncate" title="${escapeHtml(s.workspace_dir)}"><i class="fa-regular fa-folder text-[9px]"></i> ${escapeHtml(wsBase)}</span><span>•</span>` : ''}
              <span>${formattedTime}</span>
              <span>•</span>
              <span>${s.message_count || 0} 条消息</span>
            </div>
          </div>
          <button class="btn-delete-session opacity-0 group-hover:opacity-100 p-1.5 rounded-lg hover:bg-[#f28b82]/15 text-[#8e918f] hover:text-[#f28b82] transition" title="删除此会话">
            <i class="fa-regular fa-trash-can text-xs pointer-events-none"></i>
          </button>
        `;

        item.addEventListener('click', (e) => {
          if (e.target.closest('.btn-delete-session')) return;
          switchToSession(s.id);
        });

        const btnDel = item.querySelector('.btn-delete-session');
        if (btnDel) {
          btnDel.addEventListener('click', async (e) => {
            e.stopPropagation();
            if (confirm(`确定要删除会话 "${s.title}" 吗？此操作无法撤销。`)) {
              await deleteSession(s.id);
            }
          });
        }

        sessionList.appendChild(item);
      });
    } catch (e) {
      console.error('loadSessions failed', e);
    }
  }

  async function switchToSession(id) {
    try {
      const res = await fetch(`/api/sessions/detail?id=${encodeURIComponent(id)}`);
      if (!res.ok) return;
      const sess = await res.json();
      activeSessionId = sess.id;
      if (sess.workspace_dir) {
        updateWorkspaceUI(sess.workspace_dir);
        checkGitStatus();
      }
      renderMessageTree(sess.messages || []);
      if (sess.snapshots) {
        populateSnapshots(sess.snapshots);
      }
      if (sess.total_tokens !== undefined) {
        updateTelemetry({
          prompt_tokens: sess.prompt_tokens,
          completion_tokens: sess.completion_tokens,
          total_tokens: sess.total_tokens,
          cost_cny: sess.cost_cny,
          cost_usd: sess.cost_usd
        });
      }
      closeDrawer();
      loadSessions();
    } catch (e) {
      console.error('switchToSession failed', e);
    }
  }

  async function createNewSession() {
    try {
      const res = await fetch('/api/sessions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        body: JSON.stringify({ title: '新研发会话' })
      });
      if (res.ok) {
        const sess = await res.json();
        activeSessionId = sess.id;
        if (sess.workspace_dir) {
          updateWorkspaceUI(sess.workspace_dir);
          checkGitStatus();
        }
        resetChatUI();
        loadSnapshots();
        updateTelemetry({ total_tokens: 0, prompt_tokens: 0, completion_tokens: 0, cost_cny: 0, cost_usd: 0 });
        loadSessions();
        closeDrawer();
      }
    } catch (e) {
      console.error('createNewSession failed', e);
    }
  }

  async function deleteSession(id) {
    try {
      const res = await fetch(`/api/sessions/detail?id=${encodeURIComponent(id)}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        if (activeSessionId === id) {
          resetChatUI();
          activeSessionId = '';
        }
        loadSessions();
      }
    } catch (e) {
      console.error('deleteSession failed', e);
    }
  }

  if (btnDrawerNewSession) btnDrawerNewSession.addEventListener('click', createNewSession);
  if (btnResetSession) btnResetSession.addEventListener('click', createNewSession);
  if (btnPromptResetSession) btnPromptResetSession.addEventListener('click', createNewSession);

  // -------------------------------------------------------------
  // Model Profiles & Persistent Configuration
  // -------------------------------------------------------------
  async function loadModelProfiles() {
    try {
      const res = await fetch('/api/models');
      if (!res.ok) return;
      const data = await res.json();
      const profiles = data.profiles || [];
      activeProfileId = data.active_id || (profiles[0] ? profiles[0].id : 'deepseek');
      cachedProfilesData = profiles;
      if (data.thinking_definitions) {
        allThinkingDefinitions = data.thinking_definitions;
      }

      // Populate Top App Bar Model Select
      if (modelSelect) {
        modelSelect.innerHTML = '';
        profiles.forEach(p => {
          const opt = document.createElement('option');
          opt.value = p.id;
          opt.textContent = (p.name && p.name !== p.model) ? `${p.name} (${p.model})` : (p.name || p.model);
          if (p.id === activeProfileId) opt.selected = true;
          modelSelect.appendChild(opt);
        });
      }

      // Populate Settings Model Profiles List
      if (modelProfilesContainer) {
        modelProfilesContainer.innerHTML = '';
        profiles.forEach(p => {
          const isAct = p.id === activeProfileId;
          const card = document.createElement('div');
          card.className = `p-3 rounded-2xl border transition flex items-center justify-between ${
            isAct
              ? 'bg-[#8ab4f8]/10 border-[#8ab4f8]/40'
              : 'bg-[#131314] hover:bg-[#282a2c] border-[#282a2c]'
          }`;

          card.innerHTML = `
            <div class="space-y-0.5">
              <div class="flex items-center gap-2">
                <span class="font-medium text-xs text-[#e3e3e3]">${escapeHtml(p.name)}</span>
                <span class="font-mono text-[10px] px-2 py-0.5 rounded-full ${
                  isAct ? 'bg-[#8ab4f8] text-[#001d35] font-semibold' : 'bg-[#1e1f20] text-[#8e918f] border border-[#37393b]'
                }">${escapeHtml(p.model)}</span>
              </div>
              <div class="text-[10px] font-mono text-[#8e918f] max-w-sm truncate">${escapeHtml(p.base_url || '默认 API 地址')}</div>
            </div>
            <div class="flex items-center gap-2">
              ${
                isAct
                  ? '<span class="text-[11px] text-[#8ab4f8] font-medium flex items-center gap-1"><i class="fa-solid fa-check"></i> 当前使用</span>'
                  : `<button class="btn-switch-prof px-3 py-1 rounded-full bg-[#282a2c] hover:bg-[#313335] text-[#c4c7c5] hover:text-[#8ab4f8] border border-[#37393b] text-[11px] font-medium transition cursor-pointer" data-id="${p.id}">使用</button>`
              }
              ${
                profiles.length > 1
                  ? `<button class="btn-del-prof p-1.5 rounded-lg hover:bg-[#f28b82]/15 text-[#8e918f] hover:text-[#f28b82] transition cursor-pointer" data-id="${p.id}" title="删除此模型"><i class="fa-regular fa-trash-can text-xs pointer-events-none"></i></button>`
                  : ''
              }
            </div>
          `;

          const btnSwitch = card.querySelector('.btn-switch-prof');
          if (btnSwitch) {
            btnSwitch.addEventListener('click', () => switchModelProfile(p.id));
          }

          const btnDel = card.querySelector('.btn-del-prof');
          if (btnDel) {
            btnDel.addEventListener('click', () => deleteModelProfile(p.id));
          }

          modelProfilesContainer.appendChild(card);
        });
      }

      // Sync active profile inputs
      const current = profiles.find(p => p.id === activeProfileId);
      if (current) {
        if (activeProfileBadge) activeProfileBadge.textContent = current.name;
        if (cfgBaseUrl) cfgBaseUrl.value = current.base_url || '';
        if (cfgApiKey) cfgApiKey.value = current.api_key || '';
        if (cfgModelName) cfgModelName.value = current.model || '';
        if (cfgMaxSteps) cfgMaxSteps.value = current.max_steps || 15;
      }

      updateThinkingLevelUI(activeProfileId);
    } catch (e) {
      console.error('loadModelProfiles failed', e);
    }
  }

  function updateThinkingLevelUI(profileId) {
    if (!btnThinkingLevel || !thinkingLevelLabel) return;
    const current = cachedProfilesData.find(p => p.id === profileId) || cachedProfilesData[0];
    if (!current) return;

    if (thinkingModelName) {
      thinkingModelName.textContent = current.model || current.name;
    }

    const levels = current.thinking_levels || [];
    if (!levels || levels.length === 0) {
      // Current model does not support reasoning
      btnThinkingLevel.disabled = true;
      btnThinkingLevel.classList.add('opacity-40', 'cursor-not-allowed');
      btnThinkingLevel.classList.remove('hover:bg-[#282a2c]', 'hover:text-[#8ab4f8]');
      thinkingLevelLabel.textContent = '思考: 不支持';
      btnThinkingLevel.title = `当前模型 [${current.name || current.model}] 不支持深度思考/推理链能力`;
      if (thinkingLevelMenu) thinkingLevelMenu.classList.add('hidden');
      return;
    }

    // Current model supports reasoning
    btnThinkingLevel.disabled = false;
    btnThinkingLevel.classList.remove('opacity-40', 'cursor-not-allowed');
    btnThinkingLevel.classList.add('hover:bg-[#282a2c]', 'hover:text-[#8ab4f8]');

    // Validate if currentThinkingLevel is among supported levels
    if (!levels.includes(currentThinkingLevel)) {
      currentThinkingLevel = current.default_thinking_level || levels[levels.length - 1] || 'high';
      localStorage.setItem('devagent_thinking_level', currentThinkingLevel);
    }

    const def = allThinkingDefinitions[currentThinkingLevel] || { label: currentThinkingLevel, desc: '' };
    thinkingLevelLabel.textContent = def.label;
    btnThinkingLevel.title = `思考深度: ${def.label} (${def.desc || ''})\n点击切换当前模型的思考算力`;

    // Populate dropdown options
    if (thinkingLevelOptions) {
      thinkingLevelOptions.innerHTML = '';
      levels.forEach(lvl => {
        const itemDef = allThinkingDefinitions[lvl] || { label: lvl, desc: '', tokens: '' };
        const isSelected = lvl === currentThinkingLevel;
        const isDefault = lvl === current.default_thinking_level;

        const optBtn = document.createElement('button');
        optBtn.type = 'button';
        optBtn.className = `w-full text-left px-3 py-2 rounded-xl text-xs transition flex items-center justify-between cursor-pointer ${
          isSelected
            ? 'bg-[#8ab4f8]/15 text-[#8ab4f8] font-medium border border-[#8ab4f8]/30'
            : 'text-[#c4c7c5] hover:bg-[#282a2c] hover:text-[#e3e3e3] border border-transparent'
        }`;

        optBtn.innerHTML = `
          <div class="space-y-0.5 min-w-0 pr-2">
            <div class="flex items-center gap-1.5">
              <span class="font-medium text-xs">${isSelected ? '✓ ' : '  '}${escapeHtml(itemDef.label)}</span>
              ${itemDef.tokens && itemDef.tokens !== '0' ? `<span class="text-[9px] px-1.5 py-0.2 rounded-full bg-[#131314] text-[#8e918f] font-mono border border-[#37393b]">${escapeHtml(itemDef.tokens)}</span>` : ''}
              ${isDefault ? `<span class="text-[9px] text-[#81c995] font-normal">· 推荐</span>` : ''}
            </div>
            <div class="text-[10px] text-[#8e918f] leading-tight truncate">${escapeHtml(itemDef.desc)}</div>
          </div>
          ${isSelected ? '<i class="fa-solid fa-check text-[#8ab4f8] text-xs flex-shrink-0"></i>' : ''}
        `;

        optBtn.addEventListener('click', (e) => {
          e.stopPropagation();
          currentThinkingLevel = lvl;
          localStorage.setItem('devagent_thinking_level', lvl);
          updateThinkingLevelUI(profileId);
          if (thinkingLevelMenu) thinkingLevelMenu.classList.add('hidden');
        });

        thinkingLevelOptions.appendChild(optBtn);
      });
    }
  }

  // Bind Thinking Level Menu Events
  if (btnThinkingLevel && thinkingLevelMenu) {
    btnThinkingLevel.addEventListener('click', (e) => {
      e.stopPropagation();
      if (btnThinkingLevel.disabled) return;
      thinkingLevelMenu.classList.toggle('hidden');
    });

    document.addEventListener('click', (e) => {
      if (!thinkingLevelMenu.classList.contains('hidden')) {
        if (!thinkingLevelMenu.contains(e.target) && !btnThinkingLevel.contains(e.target)) {
          thinkingLevelMenu.classList.add('hidden');
        }
      }
    });
  }

  async function switchModelProfile(profileId) {
    try {
      const res = await fetch('/api/models', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        body: JSON.stringify({ action: 'switch', profile_id: profileId })
      });
      if (res.ok) {
        activeProfileId = profileId;
        await loadModelProfiles();
        updateThinkingLevelUI(profileId);
      }
    } catch (e) {
      console.error('switchModelProfile error', e);
    }
  }

  async function deleteModelProfile(profileId) {
    if (!confirm('确定要从模型库中删除此模型档案吗？')) return;
    try {
      const res = await fetch('/api/models', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        body: JSON.stringify({ action: 'delete', profile_id: profileId })
      });
      if (res.ok) {
        await loadModelProfiles();
      }
    } catch (e) {
      console.error('deleteModelProfile error', e);
    }
  }

  if (modelSelect) {
    modelSelect.addEventListener('change', () => {
      updateThinkingLevelUI(modelSelect.value);
      switchModelProfile(modelSelect.value);
    });
  }

  // Toggle Add Model Form
  if (btnToggleAddModel && addModelPanel) {
    btnToggleAddModel.addEventListener('click', () => {
      addModelPanel.classList.toggle('hidden');
    });
    if (btnCancelAddModel) {
      btnCancelAddModel.addEventListener('click', () => {
        addModelPanel.classList.add('hidden');
      });
    }
  }

  if (btnConfirmAddModel) {
    btnConfirmAddModel.addEventListener('click', async () => {
      const id = (newModelId.value || '').trim();
      const name = (newModelName.value || '').trim();
      const url = (newModelUrl.value || '').trim();
      const key = (newModelKey.value || '').trim();
      const model = (newModelCode.value || '').trim();
      const steps = parseInt(newModelSteps.value, 10) || 15;

      if (!id || !name || !url || !model) {
        alert('请完整填写模型 ID、显示名称、Base URL 和模型代码！');
        return;
      }

      try {
        const res = await fetch('/api/models', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          body: JSON.stringify({
            action: 'save',
            profile: {
              id: id,
              name: name,
              provider: 'openai_custom',
              base_url: url,
              api_key: key,
              model: model,
              max_steps: steps
            }
          })
        });

        if (res.ok) {
          addModelPanel.classList.add('hidden');
          newModelId.value = '';
          newModelName.value = '';
          newModelUrl.value = '';
          newModelKey.value = '';
          newModelCode.value = '';
          await switchModelProfile(id);
        } else {
          alert('添加模型失败: ' + (await res.text()));
        }
      } catch (e) {
        alert('添加模型异常: ' + e.message);
      }
    });
  }

  // -------------------------------------------------------------
  // Provider Hub Logic (Full Pi Agent Multi-Provider Registry)
  // -------------------------------------------------------------
  let officialProviders = [];
  let currentProvTab = 'api_key'; // 'api_key' or 'account'
  let currentProvCat = 'all';
  let provSearchQuery = '';

  const tabProvApiKey = document.getElementById('tab-prov-apikey');
  const tabProvAccount = document.getElementById('tab-prov-account');
  const provSearchInput = document.getElementById('prov-search-input');
  const provFilterBtns = document.querySelectorAll('.btn-filter-prov');
  const providerCardsGrid = document.getElementById('provider-cards-grid');

  function getProviderIcon(id) {
    const s = (id || '').toLowerCase();
    if (s.includes('deepseek')) return 'fa-solid fa-whale text-[#8ab4f8]';
    if (s.includes('anthropic') || s.includes('claude')) return 'fa-solid fa-shapes text-[#f28b82]';
    if (s.includes('openai') || s.includes('codex')) return 'fa-solid fa-brain text-[#81c995]';
    if (s.includes('gemini') || s.includes('google')) return 'fa-solid fa-sparkles text-[#a8c7fa]';
    if (s.includes('mistral')) return 'fa-solid fa-wind text-[#fdd663]';
    if (s.includes('xai')) return 'fa-solid fa-xmark text-white';
    if (s.includes('groq')) return 'fa-solid fa-bolt-lightning text-[#f28b82]';
    if (s.includes('cerebras')) return 'fa-solid fa-microchip text-[#fdd663]';
    if (s.includes('nvidia')) return 'fa-solid fa-server text-[#81c995]';
    if (s.includes('together')) return 'fa-solid fa-network-wired text-[#8ab4f8]';
    if (s.includes('moonshot') || s.includes('kimi')) return 'fa-solid fa-moon text-[#81c995]';
    if (s.includes('minimax')) return 'fa-solid fa-gem text-[#c5b3f8]';
    if (s.includes('qwen')) return 'fa-solid fa-cloud text-[#8ab4f8]';
    if (s.includes('xiaomi')) return 'fa-solid fa-mobile-screen text-[#ff6700]';
    if (s.includes('zai')) return 'fa-solid fa-bolt text-[#8ab4f8]';
    if (s.includes('copilot')) return 'fa-brands fa-github text-white';
    if (s.includes('ant-ling')) return 'fa-solid fa-shield-halved text-[#1677ff]';
    if (s.includes('meta')) return 'fa-brands fa-meta text-[#0668e1]';
    if (s.includes('openrouter')) return 'fa-solid fa-diagram-project text-[#8ab4f8]';
    if (s.includes('cloudflare')) return 'fa-brands fa-cloudflare text-[#f6821f]';
    if (s.includes('bedrock') || s.includes('amazon')) return 'fa-brands fa-aws text-[#ff9900]';
    if (s.includes('azure')) return 'fa-brands fa-microsoft text-[#00a4ef]';
    if (s.includes('huggingface')) return 'fa-regular fa-face-smile text-[#ffd21e]';
    if (s.includes('llama') || s.includes('ollama')) return 'fa-solid fa-laptop-code text-[#fdd663]';
    return 'fa-solid fa-cube text-[#8ab4f8]';
  }

  function getProviderBadge(p) {
    if (p.id === 'deepseek') return { text: '推荐', color: 'bg-[#8ab4f8]/15 text-[#8ab4f8] border-[#8ab4f8]/30' };
    if (p.type === 'account') return { text: '订阅', color: 'bg-white/10 text-white border-white/20' };
    if (p.category === 'china') return { text: '国内', color: 'bg-[#81c995]/15 text-[#81c995] border-[#81c995]/30' };
    if (p.category === 'chips') return { text: '芯片', color: 'bg-[#fdd663]/15 text-[#fdd663] border-[#fdd663]/30' };
    if (p.category === 'local') return { text: '免Key', color: 'bg-[#fdd663]/15 text-[#fdd663] border-[#fdd663]/30' };
    if (p.category === 'gateway') return { text: '网关', color: 'bg-[#c5b3f8]/15 text-[#c5b3f8] border-[#c5b3f8]/30' };
    return { text: '官方', color: 'bg-[#8ab4f8]/15 text-[#8ab4f8] border-[#8ab4f8]/30' };
  }

  async function loadOfficialProviders() {
    try {
      const res = await fetch('/api/providers');
      if (res.ok) {
        officialProviders = await res.json();
        updateProviderTabCounts();
        renderProviderCards();
      }
    } catch (e) {
      console.error('Failed to load official providers', e);
    }
  }

  function updateProviderTabCounts() {
    const apiCount = officialProviders.filter(p => p.type === 'api_key').length;
    const accCount = officialProviders.filter(p => p.type === 'account').length;
    const badgeApi = document.getElementById('badge-apikey-count');
    const badgeAcc = document.getElementById('badge-account-count');
    if (badgeApi) badgeApi.textContent = apiCount;
    if (badgeAcc) badgeAcc.textContent = accCount;
  }

  function renderProviderCards() {
    if (!providerCardsGrid) return;

    const filtered = officialProviders.filter(p => {
      if (p.type !== currentProvTab) return false;
      if (currentProvCat !== 'all' && p.category !== currentProvCat) return false;
      if (provSearchQuery) {
        const q = provSearchQuery.toLowerCase();
        const matchName = (p.name || '').toLowerCase().includes(q);
        const matchId = (p.id || '').toLowerCase().includes(q);
        const matchModel = (p.models || []).some(m => (m.id || '').toLowerCase().includes(q) || (m.name || '').toLowerCase().includes(q));
        if (!matchName && !matchId && !matchModel) return false;
      }
      return true;
    });

    if (filtered.length === 0) {
      providerCardsGrid.innerHTML = '<div class="col-span-full py-8 text-center text-xs text-[#8e918f]">未找到匹配的厂商，可尝试更换分类标签或清除搜索词。</div>';
      return;
    }

    providerCardsGrid.innerHTML = '';
    filtered.forEach(p => {
      const card = document.createElement('button');
      card.type = 'button';
      card.className = 'btn-provider-card p-2.5 rounded-2xl bg-[#191a1b] hover:bg-[#282a2c] border border-[#37393b] hover:border-[#8ab4f8]/50 text-left transition group cursor-pointer';
      card.dataset.provider = p.id;
      card.dataset.name = p.name;
      card.dataset.url = p.base_url;
      card.dataset.type = p.type;

      if (selectedProvider === p.id) {
        card.classList.add('ring-2', 'ring-[#8ab4f8]', 'bg-[#282a2c]');
      }

      const iconClass = getProviderIcon(p.id);
      const badge = getProviderBadge(p);
      const topModelNames = (p.models || []).slice(0, 2).map(m => m.id).join(' · ') || (p.models && p.models.length ? `${p.models.length} 个模型` : '官方通道');

      card.innerHTML = `
        <div class="flex items-center justify-between mb-0.5">
          <span class="font-semibold text-xs text-[#e3e3e3] group-hover:text-[#8ab4f8] flex items-center gap-1.5 truncate">
            <i class="${iconClass}"></i> <span class="truncate">${escapeHtml(p.name)}</span>
          </span>
          <span class="text-[9px] px-1.5 py-0.2 rounded border ${badge.color}">${badge.text}</span>
        </div>
        <div class="text-[10px] text-[#8e918f] truncate font-mono">${escapeHtml(topModelNames)}</div>
      `;

      card.addEventListener('click', () => {
        selectProviderCard(p);
      });

      providerCardsGrid.appendChild(card);
    });
  }

  function selectProviderCard(p) {
    selectedProvider = p.id;
    document.querySelectorAll('.btn-provider-card').forEach(c => {
      if (c.dataset.provider === p.id) {
        c.classList.add('ring-2', 'ring-[#8ab4f8]', 'bg-[#282a2c]');
      } else {
        c.classList.remove('ring-2', 'ring-[#8ab4f8]', 'bg-[#282a2c]');
      }
    });

    if (discProviderBadge) discProviderBadge.textContent = p.name;
    if (discBaseUrl) discBaseUrl.value = p.base_url || '';

    if (discApiKey) {
      discApiKey.value = '';
      if (p.type === 'account') {
        discApiKey.placeholder = p.auth_env ? `输入账户凭证 (${p.auth_env})...` : '输入账户授权 Token / Device Token...';
      } else {
        discApiKey.placeholder = p.auth_env ? `输入 API Key (${p.auth_env})...` : '输入 API Key (本地服务可留空)...';
      }
      if (cfgApiKey && cfgApiKey.value) {
        discApiKey.value = cfgApiKey.value;
      }
    }

    if (discStatusHint) discStatusHint.textContent = `已就绪，点击按钮联网探测 ${p.name} 最新模型`;

    // Render official models as instant preview
    if (p.models && p.models.length > 0) {
      renderDiscoveredModelsList(p.models, p.base_url || '', discApiKey ? discApiKey.value : '');
    } else {
      if (discModelsResult) discModelsResult.classList.add('hidden');
    }

    if (providerDiscoveryPanel) providerDiscoveryPanel.classList.remove('hidden');
  }

  function renderDiscoveredModelsList(modelsList, baseUrl, apiKey) {
    fetchedDiscoveredModels = modelsList.map(m => ({
      id: m.id,
      name: m.name || m.id,
      provider: selectedProvider,
      base_url: baseUrl || (discBaseUrl ? discBaseUrl.value.trim() : ''),
      api_key: apiKey || (discApiKey ? discApiKey.value.trim() : ''),
      description: m.description || (m.reasoning ? '深度长链推理模型 (Reasoning)' : '官方全能模型')
    }));

    if (discModelsList) {
      discModelsList.innerHTML = '';
      if (fetchedDiscoveredModels.length === 0) {
        discModelsList.innerHTML = '<div class="text-[#8e918f] p-3 text-center">未检索到可用模型，请检查 API Key 或 Base URL。</div>';
      } else {
        fetchedDiscoveredModels.forEach((m, idx) => {
          const item = document.createElement('div');
          item.className = 'flex items-start gap-2.5 p-2 rounded-xl bg-[#131314] border border-[#282a2c] hover:border-[#37393b] transition';
          item.innerHTML = `
            <input type="checkbox" id="disc-mod-${idx}" data-idx="${idx}" checked class="mt-1 rounded border-[#37393b] cursor-pointer">
            <label for="disc-mod-${idx}" class="flex-1 cursor-pointer select-text">
              <div class="flex items-center gap-2">
                <span class="font-mono font-medium text-xs text-[#8ab4f8]">${escapeHtml(m.id)}</span>
                <span class="text-[11px] text-[#e3e3e3]">${escapeHtml(m.name || m.id)}</span>
              </div>
              ${m.description ? `<div class="text-[10px] text-[#8e918f] mt-0.5">${escapeHtml(m.description)}</div>` : ''}
            </label>
          `;
          discModelsList.appendChild(item);
        });
      }
    }

    if (discFoundCount) discFoundCount.textContent = `已收录 ${fetchedDiscoveredModels.length} 个最新官方模型`;
    if (discStatusHint) discStatusHint.textContent = `共展示 ${fetchedDiscoveredModels.length} 个可用模型，勾选后点击下方按钮一键纳管落盘。`;
    if (discModelsResult) discModelsResult.classList.remove('hidden');
  }

  // Dual Tabs Switcher Handlers
  if (tabProvApiKey) {
    tabProvApiKey.addEventListener('click', () => {
      currentProvTab = 'api_key';
      tabProvApiKey.className = 'flex-1 py-2 px-3 rounded-xl text-xs font-medium flex items-center justify-center gap-2 transition bg-[#8ab4f8]/20 text-[#8ab4f8] border border-[#8ab4f8]/30 cursor-pointer shadow-sm';
      tabProvAccount.className = 'flex-1 py-2 px-3 rounded-xl text-xs font-medium flex items-center justify-center gap-2 transition text-[#8e918f] hover:text-[#e3e3e3] hover:bg-[#1e1f20] cursor-pointer';
      renderProviderCards();
    });
  }

  if (tabProvAccount) {
    tabProvAccount.addEventListener('click', () => {
      currentProvTab = 'account';
      tabProvAccount.className = 'flex-1 py-2 px-3 rounded-xl text-xs font-medium flex items-center justify-center gap-2 transition bg-[#8ab4f8]/20 text-[#8ab4f8] border border-[#8ab4f8]/30 cursor-pointer shadow-sm';
      tabProvApiKey.className = 'flex-1 py-2 px-3 rounded-xl text-xs font-medium flex items-center justify-center gap-2 transition text-[#8e918f] hover:text-[#e3e3e3] hover:bg-[#1e1f20] cursor-pointer';
      renderProviderCards();
    });
  }

  // Real-time Provider Search
  if (provSearchInput) {
    provSearchInput.addEventListener('input', () => {
      provSearchQuery = provSearchInput.value.trim();
      renderProviderCards();
    });
  }

  // Category Filter Buttons
  if (provFilterBtns) {
    provFilterBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        currentProvCat = btn.dataset.cat || 'all';
        provFilterBtns.forEach(b => {
          b.classList.remove('active', 'bg-[#8ab4f8]/20', 'text-[#8ab4f8]', 'border-[#8ab4f8]/30', 'font-medium');
          b.classList.add('bg-[#191a1b]', 'text-[#c4c7c5]', 'border-[#37393b]');
        });
        btn.classList.add('active', 'bg-[#8ab4f8]/20', 'text-[#8ab4f8]', 'border-[#8ab4f8]/30', 'font-medium');
        btn.classList.remove('bg-[#191a1b]', 'text-[#c4c7c5]', 'border-[#37393b]');
        renderProviderCards();
      });
    });
  }

  if (btnCloseDiscPanel && providerDiscoveryPanel) {
    btnCloseDiscPanel.addEventListener('click', () => {
      providerDiscoveryPanel.classList.add('hidden');
      document.querySelectorAll('.btn-provider-card').forEach(c => c.classList.remove('ring-2', 'ring-[#8ab4f8]', 'bg-[#282a2c]'));
    });
  }

  if (btnFetchProviderModels) {
    btnFetchProviderModels.addEventListener('click', async () => {
      const url = discBaseUrl ? discBaseUrl.value.trim() : '';
      const key = discApiKey ? discApiKey.value.trim() : '';

      btnFetchProviderModels.disabled = true;
      btnFetchProviderModels.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> <span>正在探测接口中...</span>';
      if (discStatusHint) discStatusHint.textContent = '正在联网请求 /models 接口探测可用模型列表...';

      try {
        const res = await fetch('/api/providers/fetch-models', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          body: JSON.stringify({
            provider: selectedProvider,
            base_url: url,
            api_key: key
          })
        });

        if (!res.ok) {
          throw new Error(await res.text());
        }

        const data = await res.json();
        const models = data.models || [];
        renderDiscoveredModelsList(models, url, key);

      } catch (err) {
        if (discStatusHint) discStatusHint.textContent = `拉取异常: ${err.message}`;
        alert('拉取模型失败: ' + err.message);
      } finally {
        btnFetchProviderModels.disabled = false;
        btnFetchProviderModels.innerHTML = '<i class="fa-solid fa-cloud-arrow-down"></i> <span>联网拉取最新可用模型</span>';
      }
    });
  }

  if (discSelectAll && discModelsList) {
    discSelectAll.addEventListener('change', () => {
      const checked = discSelectAll.checked;
      discModelsList.querySelectorAll('input[type="checkbox"]').forEach(cb => cb.checked = checked);
    });
  }

  if (btnImportDiscModels) {
    btnImportDiscModels.addEventListener('click', async () => {
      const checkboxes = discModelsList.querySelectorAll('input[type="checkbox"]:checked');
      if (checkboxes.length === 0) {
        alert('请至少勾选一个模型进行导入！');
        return;
      }

      const selectedProfiles = [];
      checkboxes.forEach(cb => {
        const idx = parseInt(cb.dataset.idx, 10);
        const m = fetchedDiscoveredModels[idx];
        if (m) {
          selectedProfiles.push({
            id: m.id,
            name: m.name || m.id,
            provider: selectedProvider,
            base_url: m.base_url || (discBaseUrl ? discBaseUrl.value.trim() : ''),
            api_key: m.api_key || (discApiKey ? discApiKey.value.trim() : ''),
            model: m.id,
            max_steps: 15
          });
        }
      });

      btnImportDiscModels.disabled = true;
      btnImportDiscModels.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> <span>正在纳管模型库并落盘...</span>';

      try {
        const res = await fetch('/api/models', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          body: JSON.stringify({
            action: 'batch_save',
            profiles: selectedProfiles
          })
        });

        if (!res.ok) {
          throw new Error(await res.text());
        }

        // Switch to the first selected profile
        if (selectedProfiles.length > 0) {
          await switchModelProfile(selectedProfiles[0].id);
        } else {
          await loadModelProfiles();
        }

        alert(`🎉 成功导入 ${selectedProfiles.length} 个最新模型！当前已自动切换为: ${selectedProfiles[0].name}`);
        if (providerDiscoveryPanel) providerDiscoveryPanel.classList.add('hidden');
      } catch (err) {
        alert('导入模型失败: ' + err.message);
      } finally {
        btnImportDiscModels.disabled = false;
        btnImportDiscModels.innerHTML = '<i class="fa-solid fa-circle-check"></i> <span>一键全部添加到模型库并使用</span>';
      }
    });
  }

  // Save Active Model Configuration (Persists to config.json)
  if (btnSaveSettings) {
    btnSaveSettings.addEventListener('click', async () => {
      const updatedProfile = {
        id: activeProfileId,
        name: activeProfileBadge ? activeProfileBadge.textContent.trim() : activeProfileId,
        provider: 'openai_compatible',
        base_url: cfgBaseUrl.value.trim(),
        api_key: cfgApiKey.value.trim(),
        model: cfgModelName.value.trim(),
        max_steps: parseInt(cfgMaxSteps.value, 10) || 15
      };

      try {
        const res = await fetch('/api/models', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          body: JSON.stringify({ action: 'save', profile: updatedProfile })
        });
        if (res.ok) {
          await loadModelProfiles();
          alert('配置已成功持久化落盘至 config.json，并在当前会话中实时生效！');
          settingsModal.classList.add('hidden');
        } else {
          alert('保存配置失败: ' + (await res.text()));
        }
      } catch (e) {
        alert('保存异常: ' + e.message);
      }
    });
  }

  // Load General App Config
  function loadConfig() {
    fetch('/api/config')
      .then(res => res.json())
      .then(data => {
        if (data.workspace) {
          updateWorkspaceUI(data.workspace);
          checkGitStatus();
        }
        if (data.approval_mode && approvalModeSelect) {
          approvalModeSelect.value = data.approval_mode;
        }
      })
      .catch(err => console.log('Config load fallback to defaults', err));
  }

  // -------------------------------------------------------------
  // WebSocket Communication & Streaming Lifecycle
  // -------------------------------------------------------------
  let pingInterval = null;

  function connectWS() {
    if (pingInterval) {
      clearInterval(pingInterval);
      pingInterval = null;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      updateStatus('idle', '已就绪');
      checkGitStatus();
      loadSnapshots();
      loadSessions();
      loadModelProfiles();

      // Start keep-alive heartbeat ping every 15s to keep connection alive
      pingInterval = setInterval(() => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          try {
            ws.send(JSON.stringify({ type: 'ping' }));
          } catch (e) {
            console.warn('Ping error:', e);
          }
        }
      }, 15000);
    };

    ws.onclose = () => {
      if (pingInterval) {
        clearInterval(pingInterval);
        pingInterval = null;
      }
      updateStatus('offline', '连接断开，重试中...');
      setTimeout(connectWS, 2000);
    };

    ws.onerror = (err) => {
      console.error('WS Error:', err);
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        handleServerMessage(msg);
      } catch (e) {
        console.error('Failed to parse message:', event.data);
      }
    };
  }

  function updateStatus(state, text) {
    statusText.textContent = text;
    statusDot.className = 'w-2 h-2 rounded-full';
    if (state === 'idle') {
      statusDot.classList.add('bg-[#81c995]');
      btnSend.disabled = false;
      btnStop.classList.add('hidden');
    } else if (state === 'thinking' || state === 'executing') {
      statusDot.classList.add('bg-[#8ab4f8]', 'animate-pulse');
      btnSend.disabled = true;
      btnStop.classList.remove('hidden');
    } else if (state === 'awaiting_approval') {
      statusDot.classList.add('bg-[#fdd663]', 'animate-ping');
      btnSend.disabled = true;
      btnStop.classList.remove('hidden');
    } else {
      statusDot.classList.add('bg-[#f28b82]');
    }
  }

  let activeMessageCard = null;
  let currentBlockType = 'none'; // 'none' | 'thought' | 'tool' | 'content'
  let currentThoughtBox = null;
  let currentThoughtContent = null;
  let currentContentDiv = null;
  let currentContentBuffer = '';

  // User Scroll Intent Tracking (Solves mouse wheel freeze & scroll hijacking)
  let userScrolledUp = false;
  if (messagesContainer) {
    messagesContainer.addEventListener('scroll', () => {
      const threshold = 60; // px tolerance
      const isNearBottom = messagesContainer.scrollHeight - messagesContainer.scrollTop - messagesContainer.clientHeight <= threshold;
      userScrolledUp = !isNearBottom;
    });
  }

  // Smooth auto-scroll respecting user scroll intent (prevents wheel hijacking)
  function scrollToBottom(force = false) {
    if (!messagesContainer) return;
    if (force) {
      userScrolledUp = false;
    }
    if (!userScrolledUp) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  // Robust zero-dependency Markdown parser fallback (guarantees content renders even when marked is blocked/unavailable)
  function fallbackMarkdownParse(text) {
    if (!text) return '';
    const codeBlocks = [];
    let processed = text.replace(/```([a-zA-Z0-9_-]*)\r?\n([\s\S]*?)```/g, (match, lang, code) => {
      const idx = codeBlocks.length;
      const langBadge = lang ? `<div class="text-[10px] text-[#8ab4f8] font-mono px-3 py-1 bg-[#282a2c] border-b border-[#37393b] rounded-t-xl">${escapeHtml(lang)}</div>` : '';
      codeBlocks.push(`
        <div class="my-2.5 rounded-xl border border-[#37393b] overflow-hidden bg-[#1e1f20]">
          ${langBadge}
          <pre class="p-3 text-xs font-mono text-[#e3e3e3] overflow-x-auto leading-relaxed"><code>${escapeHtml(code.trimEnd())}</code></pre>
        </div>
      `);
      return `@@CODE_BLOCK_${idx}@@`;
    });

    processed = escapeHtml(processed);

    // Headers
    processed = processed.replace(/^### (.*$)/gim, '<h3 class="text-base font-semibold text-[#8ab4f8] mt-3 mb-1.5">$1</h3>');
    processed = processed.replace(/^## (.*$)/gim, '<h2 class="text-lg font-bold text-[#e3e3e3] mt-4 mb-2">$1</h2>');
    processed = processed.replace(/^# (.*$)/gim, '<h1 class="text-xl font-bold text-[#e3e3e3] mt-4 mb-2">$1</h1>');

    // Bold
    processed = processed.replace(/\*\*(.*?)\*\*/g, '<strong class="font-semibold text-[#e3e3e3]">$1</strong>');
    processed = processed.replace(/__(.*?)__/g, '<strong class="font-semibold text-[#e3e3e3]">$1</strong>');

    // Inline code
    processed = processed.replace(/`([^`]+)`/g, '<code class="bg-[#282a2c] text-[#8ab4f8] px-1.5 py-0.5 rounded text-xs font-mono">$1</code>');

    // Blockquote
    processed = processed.replace(/^&gt;(?: |&nbsp;)?(.*$)/gim, '<blockquote class="border-l-2 border-[#8ab4f8] pl-3 my-2 text-[#c4c7c5] italic">$1</blockquote>');

    // Horizontal rule
    processed = processed.replace(/^(?:---|___|\*\*\*)$/gim, '<hr class="border-[#282a2c] my-3">');

    // Lists
    processed = processed.replace(/^[\*\-] (.*$)/gim, '<div class="flex items-start gap-2 my-1"><span class="text-[#8ab4f8]">•</span><span>$1</span></div>');
    processed = processed.replace(/^(\d+)\. (.*$)/gim, '<div class="flex items-start gap-2 my-1"><span class="text-[#8ab4f8] font-mono text-xs">$1.</span><span>$2</span></div>');

    // Links
    processed = processed.replace(/\[(.*?)\]\((https?:\/\/[^\s)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-[#8ab4f8] hover:underline">$1</a>');

    // Newlines to <br>
    processed = processed.replace(/\n/g, '<br>');

    // Restore protected code blocks
    codeBlocks.forEach((blockHtml, idx) => {
      processed = processed.replace(`@@CODE_BLOCK_${idx}@@`, blockHtml);
    });

    return processed;
  }

  function safeRenderMarkdown(text) {
    if (!text) return '';
    if (typeof marked !== 'undefined' && typeof marked.parse === 'function') {
      try {
        return marked.parse(text);
      } catch (err) {
        console.warn('marked.parse error, using fallback:', err);
      }
    }
    return fallbackMarkdownParse(text);
  }

  // Throttled Markdown Rendering with requestAnimationFrame (Solves CPU lockup during streaming)
  let renderFrameRequested = false;
  let pendingRenderBuffer = '';
  let pendingRenderDiv = null;

  function scheduleContentRender(div, text) {
    pendingRenderDiv = div;
    pendingRenderBuffer = text;
    if (!renderFrameRequested) {
      renderFrameRequested = true;
      requestAnimationFrame(() => {
        renderFrameRequested = false;
        if (pendingRenderDiv && pendingRenderBuffer) {
          try {
            pendingRenderDiv.innerHTML = safeRenderMarkdown(pendingRenderBuffer);
          } catch (e) {
            pendingRenderDiv.textContent = pendingRenderBuffer;
          }
          scrollToBottom();
        }
      });
    }
  }

  function flushContentRender() {
    renderFrameRequested = false;
    if (pendingRenderDiv && pendingRenderBuffer) {
      try {
        pendingRenderDiv.innerHTML = safeRenderMarkdown(pendingRenderBuffer);
      } catch (e) {
        pendingRenderDiv.textContent = pendingRenderBuffer;
      }
      scrollToBottom();
    }
    pendingRenderDiv = null;
    pendingRenderBuffer = '';
  }

  function handleServerMessage(msg) {
    switch (msg.type) {
      case 'pong':
        break;

      case 'status':
        updateStatus(msg.state, msg.text || msg.state);
        break;

      case 'thought':
        appendThoughtChunk(msg.content);
        break;

      case 'team_role':
        setAssistantTeamRole(msg.role, msg.name, msg.icon);
        break;

      case 'content':
        appendAssistantContent(msg.content);
        break;

      case 'tool_call':
        renderToolCall(msg.tool, msg.args, msg.call_id);
        break;

      case 'cmd_stream':
        appendTerminalStream(msg.call_id, msg.chunk);
        break;

      case 'tool_result':
        renderToolResult(msg.call_id, msg.output, msg.error);
        break;

      case 'telemetry':
        updateTelemetry(msg);
        break;

      case 'snapshot_created':
        loadSnapshots();
        loadSessions();
        break;

      case 'session_rollback':
        // Fix snapshot rollback failure: completely re-render message tree
        if (msg.messages && Array.isArray(msg.messages)) {
          renderMessageTree(msg.messages);
        }
        loadSnapshots();
        loadSessions();
        break;

      case 'session_loaded':
        if (msg.session) {
          activeSessionId = msg.session.id;
          if (msg.session.workspace_dir) {
            updateWorkspaceUI(msg.session.workspace_dir);
            checkGitStatus();
          }
          if (msg.session.messages) {
            renderMessageTree(msg.session.messages);
          }
          if (msg.session.snapshots) populateSnapshots(msg.session.snapshots);
          if (msg.session.total_tokens !== undefined) {
            updateTelemetry({
              prompt_tokens: msg.session.prompt_tokens,
              completion_tokens: msg.session.completion_tokens,
              total_tokens: msg.session.total_tokens,
              cost_cny: msg.session.cost_cny,
              cost_usd: msg.session.cost_usd
            });
          }
        }
        loadSessions();
        break;

      case 'model_switched':
        loadModelProfiles();
        break;

      case 'awaiting_approval':
        showApprovalPrompt(msg.token, msg.tool, msg.desc, msg.payload);
        break;

      case 'workspace_updated':
        if (msg.workspace) {
          updateWorkspaceUI(msg.workspace);
          checkGitStatus();
        }
        break;

      case 'approval_mode_updated':
        if (msg.approval_mode && approvalModeSelect) {
          approvalModeSelect.value = msg.approval_mode;
        }
        break;

      case 'session_reset':
        resetChatUI();
        if (msg.workspace) {
          updateWorkspaceUI(msg.workspace);
          checkGitStatus();
        }
        updateTelemetry({ total_tokens: 0, prompt_tokens: 0, completion_tokens: 0, cost_cny: 0, cost_usd: 0 });
        updateStatus('idle', '会话已重置');
        loadSnapshots();
        loadSessions();
        break;

      case 'done':
        flushContentRender();
        updateStatus('idle', '任务完成');
        finalizeThoughtAccordion();
        // Automatically collapse tool calls upon task completion so user sees final answer without endless scrolling
        document.querySelectorAll('.tool-details-container').forEach(el => el.classList.add('hidden'));
        document.querySelectorAll('.tool-header .chevron-icon').forEach(el => el.classList.remove('rotate-180'));
        approvalCard.classList.add('hidden');
        activeMessageCard = null;
        currentBlockType = 'none';
        currentThoughtBox = null;
        currentThoughtContent = null;
        currentContentDiv = null;
        currentContentBuffer = '';
        checkGitStatus();
        loadSnapshots();
        loadSessions();
        break;

      case 'error':
        appendErrorMessage(msg.error);
        updateStatus('idle', '出现错误');
        approvalCard.classList.add('hidden');
        break;
    }
  }

  function hideWelcomeHero() {
    if (welcomeHero) {
      welcomeHero.classList.add('hidden');
    }
  }

  function appendUserMessage(text) {
    hideWelcomeHero();
    const el = document.createElement('div');
    el.className = 'flex items-start justify-end max-w-3xl mx-auto my-3';
    el.innerHTML = `
      <div class="bg-[#282a2c] text-[#e3e3e3] rounded-[24px] px-5 py-3 max-w-[85%] text-sm leading-relaxed border border-[#37393b]/70 shadow-sm whitespace-pre-wrap select-text">
        ${escapeHtml(text)}
      </div>
    `;
    messagesContainer.appendChild(el);
    scrollToBottom(true);
  }

  function ensureAssistantCard() {
    hideWelcomeHero();
    if (!activeMessageCard) {
      activeMessageCard = document.createElement('div');
      activeMessageCard.className = 'flex items-start space-x-3.5 max-w-3xl mx-auto my-4';
      activeMessageCard.innerHTML = `
        <div class="w-8 h-8 rounded-full gemini-gradient flex items-center justify-center flex-shrink-0 text-white text-xs mt-0.5 shadow-md shadow-indigo-500/10">
          <i class="fa-solid fa-sparkles"></i>
        </div>
        <div class="flex-1 space-y-3 msg-body text-sm text-[#e3e3e3] min-w-0 select-text">
          <div class="team-badge-container flex items-center gap-2 mb-1.5 hidden"></div>
        </div>
      `;
      messagesContainer.appendChild(activeMessageCard);
      currentBlockType = 'none';
      currentThoughtBox = null;
      currentThoughtContent = null;
      currentContentDiv = null;
      currentContentBuffer = '';
    }
    return activeMessageCard.querySelector('.msg-body');
  }

  // -------------------------------------------------------------
  // Multi-Agent Team Orchestration Badge
  // -------------------------------------------------------------
  function setAssistantTeamRole(role, name, icon) {
    if (!activeMessageCard) {
      ensureAssistantCard();
    }
    const container = activeMessageCard.querySelector('.team-badge-container');
    if (!container) return;

    container.classList.remove('hidden');
    let colorClass = 'bg-[#8ab4f8]/15 text-[#8ab4f8] border-[#8ab4f8]/30';
    let rolePrefix = '🏛️ 架构设计';

    if (role === 'engineer') {
      colorClass = 'bg-[#81c995]/15 text-[#81c995] border-[#81c995]/30';
      rolePrefix = '💻 全栈研发';
    } else if (role === 'qa') {
      colorClass = 'bg-[#fdd663]/15 text-[#fdd663] border-[#fdd663]/30';
      rolePrefix = '🛡️ 测试审计';
    } else if (role === 'solo') {
      colorClass = 'bg-[#c5b3f8]/15 text-[#c5b3f8] border-[#c5b3f8]/30';
      rolePrefix = '⚡ 单兵极简';
    }

    container.innerHTML = `
      <span class="px-2.5 py-0.5 rounded-full text-[10px] font-medium border flex items-center gap-1.5 ${colorClass}">
        <i class="fa-solid ${icon || 'fa-user-gear'} text-[10px]"></i>
        <span>虚拟工程团队 • ${escapeHtml(name || rolePrefix)}</span>
      </span>
    `;
  }

  function appendThoughtChunk(text) {
    const body = ensureAssistantCard();
    if (currentBlockType !== 'thought' || !currentThoughtBox) {
      currentThoughtBox = document.createElement('div');
      currentThoughtBox.className = 'thought-box bg-[#191a1b] border border-[#282a2c] rounded-2xl p-3 text-xs text-[#c4c7c5] my-2 transition-all shadow-sm';
      currentThoughtBox.innerHTML = `
        <div class="thought-header flex items-center justify-between cursor-pointer text-[#8ab4f8] font-medium py-0.5 select-none hover:text-[#a8c7fa] transition">
          <span class="flex items-center gap-2">
            <i class="fa-solid fa-sparkles text-xs animate-subtle"></i>
            <span class="thought-title">深度推导与规划中...</span>
          </span>
          <div class="flex items-center gap-1.5 text-[#8e918f]">
            <span class="thought-count text-[10px] font-mono"></span>
            <i class="fa-solid fa-chevron-down text-[10px] transform transition-transform duration-200 chevron-icon"></i>
          </div>
        </div>
        <div class="thought-content mt-2.5 pt-2 border-t border-[#282a2c] text-[#c4c7c5] font-mono text-[11px] whitespace-pre-wrap leading-relaxed max-h-56 overflow-y-auto"></div>
      `;

      const header = currentThoughtBox.querySelector('.thought-header');
      const content = currentThoughtBox.querySelector('.thought-content');
      const chevron = currentThoughtBox.querySelector('.chevron-icon');
      header.addEventListener('click', (e) => {
        e.stopPropagation();
        content.classList.toggle('hidden');
        chevron.classList.toggle('rotate-180');
      });

      body.appendChild(currentThoughtBox);
      currentThoughtContent = content;
      currentBlockType = 'thought';
      currentContentDiv = null;
      currentContentBuffer = '';
    }

    currentThoughtContent.textContent += text;
    currentThoughtContent.scrollTop = currentThoughtContent.scrollHeight;

    const countEl = currentThoughtBox.querySelector('.thought-count');
    if (countEl) {
      countEl.textContent = `${currentThoughtContent.textContent.length} 字`;
    }

    scrollToBottom();
  }

  function collapseActiveThought() {
    if (currentThoughtBox) {
      const title = currentThoughtBox.querySelector('.thought-title');
      if (title) title.textContent = '思考过程 (已折叠 · 点击展开)';
      const content = currentThoughtBox.querySelector('.thought-content');
      const chevron = currentThoughtBox.querySelector('.chevron-icon');
      if (content && !content.classList.contains('hidden')) {
        content.classList.add('hidden');
        if (chevron) chevron.classList.remove('rotate-180');
      }
    }
  }

  function finalizeThoughtAccordion() {
    flushContentRender();
    if (!activeMessageCard) return;
    activeMessageCard.querySelectorAll('.thought-box').forEach(box => {
      const title = box.querySelector('.thought-title');
      if (title) title.textContent = '思考过程 (已折叠 · 点击展开)';
      const content = box.querySelector('.thought-content');
      const chevron = box.querySelector('.chevron-icon');
      if (content && !content.classList.contains('hidden')) {
        content.classList.add('hidden');
        if (chevron) chevron.classList.remove('rotate-180');
      }
    });
  }

  function collapseActiveTools() {
    if (!activeMessageCard) return;
    activeMessageCard.querySelectorAll('.tool-details-container').forEach(el => el.classList.add('hidden'));
    activeMessageCard.querySelectorAll('.tool-header .chevron-icon').forEach(el => el.classList.remove('rotate-180'));
  }

  function appendAssistantContent(chunk) {
    const body = ensureAssistantCard();
    collapseActiveThought();

    if (currentBlockType !== 'content' || !currentContentDiv) {
      collapseActiveTools();
      currentContentDiv = document.createElement('div');
      currentContentDiv.className = 'prose prose-invert prose-sm max-w-none text-[#e3e3e3] leading-relaxed my-2 select-text';
      body.appendChild(currentContentDiv);
      currentContentBuffer = '';
      currentBlockType = 'content';
    }

    currentContentBuffer += chunk;
    scheduleContentRender(currentContentDiv, currentContentBuffer);
  }

  function createToolAccordionCard(tool, args, output, isCollapsed) {
    const card = document.createElement('div');
    card.className = 'border border-[#37393b] rounded-2xl bg-[#1e1f20] p-3 text-xs font-mono my-2.5 shadow-sm transition-all';

    let formattedArgs = '';
    try {
      formattedArgs = typeof args === 'string' ? JSON.stringify(JSON.parse(args), null, 2) : JSON.stringify(args, null, 2);
    } catch (e) {
      formattedArgs = String(args || '');
    }
    const shortArgs = formattedArgs.length > 40 ? formattedArgs.slice(0, 40).replace(/\n/g, ' ') + '...' : formattedArgs.replace(/\n/g, ' ');

    card.innerHTML = `
      <div class="tool-header flex items-center justify-between text-[#c4c7c5] font-medium cursor-pointer select-none hover:text-[#e3e3e3] transition">
        <span class="flex items-center gap-2">
          <div class="w-5 h-5 rounded-md bg-[#8ab4f8]/10 text-[#8ab4f8] flex items-center justify-center text-[10px]">
            <i class="fa-solid fa-screwdriver-wrench"></i>
          </div>
          <span class="text-xs">调用工具: <strong class="text-[#8ab4f8] font-semibold">${escapeHtml(tool)}</strong></span>
          <span class="text-[10px] text-[#8e918f] font-normal truncate max-w-[200px] hidden sm:inline">(${escapeHtml(shortArgs)})</span>
        </span>
        <div class="flex items-center gap-2">
          <span class="text-[11px] px-2 py-0.5 rounded-full ${output ? 'bg-[#81c995]/15 text-[#81c995] border border-[#81c995]/30' : 'bg-[#8ab4f8]/10 text-[#8ab4f8] border border-[#8ab4f8]/20'} tool-status">
            ${output ? '已完成' : '已记录'}
          </span>
          <i class="fa-solid fa-chevron-down text-[10px] transform transition-transform duration-200 chevron-icon ${isCollapsed ? '' : 'rotate-180'}"></i>
        </div>
      </div>
      <div class="tool-details-container ${isCollapsed ? 'hidden' : ''} mt-2.5 space-y-2">
        <div class="text-[#8e918f] bg-[#131314] p-2.5 rounded-xl border border-[#282a2c] overflow-x-auto text-[11px] leading-relaxed">
          <div class="text-[10px] text-[#8ab4f8] mb-1 font-medium"><i class="fa-solid fa-code"></i> 参数 (Arguments):</div>
          <pre>${escapeHtml(formattedArgs)}</pre>
        </div>
        ${output ? `
          <div class="tool-output-area text-[11px]">
            <div class="text-[#81c995] mb-1 flex items-center gap-1"><i class="fa-solid fa-circle-check"></i> 执行结果:</div>
            <pre class="bg-[#131314] p-2.5 rounded-xl border border-[#282a2c] text-[#c4c7c5] overflow-x-auto max-h-48 whitespace-pre-wrap">${escapeHtml(output)}</pre>
          </div>
        ` : ''}
      </div>
    `;

    const header = card.querySelector('.tool-header');
    const details = card.querySelector('.tool-details-container');
    const chevron = card.querySelector('.chevron-icon');
    header.addEventListener('click', (e) => {
      e.stopPropagation();
      details.classList.toggle('hidden');
      chevron.classList.toggle('rotate-180');
    });

    return card;
  }

  function renderToolCall(tool, args, callId) {
    const body = ensureAssistantCard();
    collapseActiveThought();
    flushContentRender();

    const toolEl = document.createElement('div');
    toolEl.id = `tool-${callId}`;
    toolEl.className = 'border border-[#37393b] rounded-2xl bg-[#1e1f20] p-3 text-xs font-mono my-2.5 shadow-sm transition-all';

    let formattedArgs = '';
    try {
      formattedArgs = typeof args === 'string' ? JSON.stringify(JSON.parse(args), null, 2) : JSON.stringify(args, null, 2);
    } catch (e) {
      formattedArgs = String(args || '');
    }
    const shortArgs = formattedArgs.length > 40 ? formattedArgs.slice(0, 40).replace(/\n/g, ' ') + '...' : formattedArgs.replace(/\n/g, ' ');

    toolEl.innerHTML = `
      <div class="tool-header flex items-center justify-between text-[#c4c7c5] font-medium cursor-pointer select-none hover:text-[#e3e3e3] transition">
        <span class="flex items-center gap-2">
          <div class="w-5 h-5 rounded-md bg-[#8ab4f8]/10 text-[#8ab4f8] flex items-center justify-center text-[10px]">
            <i class="fa-solid fa-screwdriver-wrench"></i>
          </div>
          <span class="text-xs">调用工具: <strong class="text-[#8ab4f8] font-semibold">${escapeHtml(tool)}</strong></span>
          <span class="text-[10px] text-[#8e918f] font-normal truncate max-w-[200px] hidden sm:inline">(${escapeHtml(shortArgs)})</span>
        </span>
        <div class="flex items-center gap-2">
          <span class="text-[11px] px-2 py-0.5 rounded-full bg-[#8ab4f8]/10 text-[#8ab4f8] border border-[#8ab4f8]/20 animate-subtle tool-status">执行中...</span>
          <i class="fa-solid fa-chevron-down text-[10px] transform transition-transform duration-200 chevron-icon rotate-180"></i>
        </div>
      </div>
      <div class="tool-details-container mt-2.5 space-y-2">
        <div class="text-[#8e918f] bg-[#131314] p-2.5 rounded-xl border border-[#282a2c] overflow-x-auto text-[11px] leading-relaxed">
          <div class="text-[10px] text-[#8ab4f8] mb-1 font-medium"><i class="fa-solid fa-code"></i> 参数 (Arguments):</div>
          <pre>${escapeHtml(formattedArgs)}</pre>
        </div>
        <div class="tool-stream-area hidden"></div>
        <div class="tool-output-area hidden text-[11px]"></div>
      </div>
    `;

    const header = toolEl.querySelector('.tool-header');
    const details = toolEl.querySelector('.tool-details-container');
    const chevron = toolEl.querySelector('.chevron-icon');
    header.addEventListener('click', (e) => {
      e.stopPropagation();
      details.classList.toggle('hidden');
      chevron.classList.toggle('rotate-180');
    });

    body.appendChild(toolEl);
    currentBlockType = 'tool';
    currentThoughtBox = null;
    currentThoughtContent = null;
    currentContentDiv = null;
    currentContentBuffer = '';
    scrollToBottom();
  }

  function appendTerminalStream(callId, chunk) {
    const toolEl = document.getElementById(`tool-${callId}`);
    if (!toolEl) return;
    const streamArea = toolEl.querySelector('.tool-stream-area');
    if (!streamArea) return;
    streamArea.classList.remove('hidden');

    let streamPre = streamArea.querySelector('pre');
    if (!streamPre) {
      streamArea.innerHTML = `
        <div class="flex items-center justify-between text-[10px] text-[#81c995] bg-[#0c0d0e] px-3 py-1.5 rounded-t-xl border border-[#282a2c] border-b-0">
          <span class="flex items-center gap-1.5"><i class="fa-solid fa-terminal animate-pulse"></i> 终端输出流 (Console Output)</span>
          <span class="px-2 py-0.5 rounded-full bg-[#81c995]/15 text-[#81c995] text-[9px] border border-[#81c995]/20 font-mono">Live</span>
        </div>
        <pre class="bg-[#0c0d0e] text-[#81c995] p-3 rounded-b-xl border border-[#282a2c] text-[11px] font-mono overflow-x-auto max-h-48 whitespace-pre-wrap leading-relaxed shadow-inner"></pre>
      `;
      streamPre = streamArea.querySelector('pre');
    }
    streamPre.textContent += chunk;
    streamPre.scrollTop = streamPre.scrollHeight;
    scrollToBottom();
  }

  function renderToolResult(callId, output, error) {
    const toolEl = document.getElementById(`tool-${callId}`);
    if (!toolEl) return;
    const status = toolEl.querySelector('.tool-status');
    const outArea = toolEl.querySelector('.tool-output-area');
    outArea.classList.remove('hidden');

    if (error) {
      status.textContent = '失败';
      status.className = 'text-[11px] px-2 py-0.5 rounded-full bg-[#f28b82]/15 text-[#f28b82] border border-[#f28b82]/30 tool-status';
      outArea.innerHTML = `<div class="text-[#f28b82]"><i class="fa-solid fa-circle-xmark mr-1"></i> 执行错误: ${escapeHtml(error)}</div>`;
    } else {
      status.textContent = '已完成';
      status.className = 'text-[11px] px-2 py-0.5 rounded-full bg-[#81c995]/15 text-[#81c995] border border-[#81c995]/30 tool-status';
      outArea.innerHTML = `
        <div class="text-[#81c995] mb-1 flex items-center gap-1"><i class="fa-solid fa-circle-check"></i> 执行完毕</div>
        <pre class="bg-[#131314] p-2.5 rounded-xl border border-[#282a2c] text-[#c4c7c5] overflow-x-auto max-h-48 whitespace-pre-wrap">${escapeHtml(output || '(无输出)')}</pre>
      `;
    }
    scrollToBottom();
  }

  function updateTelemetry(data) {
    if (!data) return;
    if (telemetryTokens && data.total_tokens !== undefined) {
      telemetryTokens.textContent = `${data.total_tokens.toLocaleString()} Tok`;
    }
    if (telemetryCost) {
      let cny = 0;
      if (data.cost_cny !== undefined && data.cost_cny !== null) {
        cny = Number(data.cost_cny);
      } else if (data.cost_usd !== undefined && data.cost_usd !== null) {
        cny = Number(data.cost_usd) * 7.25;
      } else if (data.estimated_cost_usd !== undefined && data.estimated_cost_usd !== null) {
        cny = Number(data.estimated_cost_usd) * 7.25;
      }
      if (cny === 0 && (data.total_tokens || 0) > 0) {
        cny = 0.0001; // Avoid displaying ¥0.0000 when tokens were consumed
      }
      telemetryCost.textContent = `¥${cny.toFixed(4)}`;
    }
    if (sessionTokenCount && data.total_tokens !== undefined) {
      const p = data.prompt_tokens || 0;
      const c = data.completion_tokens || 0;
      sessionTokenCount.textContent = `${data.total_tokens.toLocaleString()} Tokens (Prompt: ${p.toLocaleString()}, Compl: ${c.toLocaleString()})`;
    }
  }

  function appendErrorMessage(err) {
    const body = ensureAssistantCard();
    const errEl = document.createElement('div');
    errEl.className = 'p-3.5 bg-[#f28b82]/10 border border-[#f28b82]/30 rounded-2xl text-xs text-[#f28b82] flex items-start gap-2.5 my-2';
    errEl.innerHTML = `<i class="fa-solid fa-circle-exclamation text-[#f28b82] mt-0.5 text-sm"></i> <span>${escapeHtml(err)}</span>`;
    body.appendChild(errEl);
    scrollToBottom();
  }

  function resetChatUI() {
    messagesContainer.innerHTML = '';
    if (welcomeHero) {
      welcomeHero.classList.remove('hidden');
      messagesContainer.appendChild(welcomeHero);
    }
    activeMessageCard = null;
    currentBlockType = 'none';
    currentThoughtBox = null;
    currentThoughtContent = null;
    currentContentDiv = null;
    currentContentBuffer = '';
    userScrolledUp = false;
  }

  // -------------------------------------------------------------
  // Full Message Tree Renderer (Solves Snapshot Rollback Visual Bug)
  // -------------------------------------------------------------
  function renderMessageTree(messages) {
    resetChatUI();
    if (!messages || !Array.isArray(messages) || messages.length === 0) {
      return;
    }
    hideWelcomeHero();

    const toolResults = {};
    for (const msg of messages) {
      if (msg.role === 'tool' && msg.tool_call_id) {
        toolResults[msg.tool_call_id] = msg.content;
      }
    }

    for (const msg of messages) {
      try {
        if (msg.role === 'user') {
          activeMessageCard = null;
          appendUserMessage(msg.content);
          activeMessageCard = null;
        } else if (msg.role === 'assistant') {
          const body = ensureAssistantCard();

          // 1. Restore historical Thought accordion
          if (msg.thought && msg.thought.trim()) {
            const thoughtBox = document.createElement('div');
            thoughtBox.className = 'thought-box bg-[#191a1b] border border-[#282a2c] rounded-2xl p-3 text-xs text-[#c4c7c5] my-2 transition-all shadow-sm';
            thoughtBox.innerHTML = `
              <div class="thought-header flex items-center justify-between cursor-pointer text-[#8ab4f8] font-medium py-0.5 select-none hover:text-[#a8c7fa] transition">
                <span class="flex items-center gap-2">
                  <i class="fa-solid fa-sparkles text-xs"></i>
                  <span class="thought-title">思考过程 (已折叠 · 点击展开)</span>
                </span>
                <div class="flex items-center gap-1.5 text-[#8e918f]">
                  <span class="text-[10px] font-mono">${msg.thought.length} 字</span>
                  <i class="fa-solid fa-chevron-down text-[10px] transform transition-transform duration-200 chevron-icon"></i>
                </div>
              </div>
              <div class="thought-content mt-2.5 pt-2 border-t border-[#282a2c] text-[#c4c7c5] font-mono text-[11px] whitespace-pre-wrap leading-relaxed max-h-56 overflow-y-auto hidden"></div>
            `;
            thoughtBox.querySelector('.thought-content').textContent = msg.thought;
            const header = thoughtBox.querySelector('.thought-header');
            const content = thoughtBox.querySelector('.thought-content');
            const chevron = thoughtBox.querySelector('.chevron-icon');
            header.addEventListener('click', (e) => {
              e.stopPropagation();
              content.classList.toggle('hidden');
              chevron.classList.toggle('rotate-180');
            });
            body.appendChild(thoughtBox);
          }

          // 2. Restore historical Tool Calls in collapsed accordions
          if (msg.tool_calls && Array.isArray(msg.tool_calls)) {
            for (const tc of msg.tool_calls) {
              const toolName = tc.function ? tc.function.name : 'tool';
              const args = tc.function ? tc.function.arguments : '';
              const output = toolResults[tc.id] || '';
              const toolCard = createToolAccordionCard(toolName, args, output, true);
              body.appendChild(toolCard);
            }
          }

          // 3. Render Assistant Content
          if (msg.content && msg.content.trim()) {
            const contentDiv = document.createElement('div');
            contentDiv.className = 'prose prose-invert prose-sm max-w-none text-[#e3e3e3] leading-relaxed my-2 select-text';
            contentDiv.innerHTML = safeRenderMarkdown(msg.content);
            body.appendChild(contentDiv);
          }
        }
      } catch (err) {
        console.error('Error rendering message item:', err, msg);
      }
    }
    activeMessageCard = null;
    scrollToBottom(true);
  }

  function showApprovalPrompt(token, tool, desc, payload) {
    pendingApprovalToken = token;
    approvalToolBadge.textContent = tool;
    approvalDesc.textContent = desc || `Agent 正在请求执行 ${tool}：`;

    if (tool === 'file_edit' && payload && typeof payload === 'object') {
      const path = payload.path || '';
      const target = payload.target_content || '';
      const repl = payload.replacement_content || '';
      approvalPayload.innerHTML = `
        <div class="text-[#e3e3e3] font-medium mb-1.5">目标文件: <span class="text-[#8ab4f8] font-mono">${escapeHtml(path)}</span></div>
        <div class="text-[#f28b82] font-semibold text-[11px] mt-2 flex items-center gap-1.5"><i class="fa-solid fa-minus"></i> 待移除/替换原代码片段:</div>
        <pre class="bg-[#1e1f20] border border-[#f28b82]/30 p-2.5 rounded-xl text-[#f28b82] text-[11px] font-mono overflow-x-auto whitespace-pre-wrap max-h-36">${escapeHtml(target)}</pre>
        <div class="text-[#81c995] font-semibold text-[11px] mt-2 flex items-center gap-1.5"><i class="fa-solid fa-plus"></i> 替换后的新代码片段:</div>
        <pre class="bg-[#1e1f20] border border-[#81c995]/30 p-2.5 rounded-xl text-[#81c995] text-[11px] font-mono overflow-x-auto whitespace-pre-wrap max-h-36">${escapeHtml(repl)}</pre>
      `;
    } else {
      approvalPayload.textContent = typeof payload === 'string' ? payload : JSON.stringify(payload, null, 2);
    }

    approvalRejectReason.value = '';
    approvalCard.classList.remove('hidden');
    updateStatus('awaiting_approval', '等待人工审核');
    scrollToBottom();
  }

  btnConfirmApproval.addEventListener('click', () => {
    if (!pendingApprovalToken || !ws) return;
    ws.send(JSON.stringify({
      type: 'approval_response',
      token: pendingApprovalToken,
      approved: true
    }));
    approvalCard.classList.add('hidden');
    pendingApprovalToken = null;
    updateStatus('executing', '操作已批准，正在执行...');
  });

  btnRejectApproval.addEventListener('click', () => {
    if (!pendingApprovalToken || !ws) return;
    const reason = approvalRejectReason.value.trim();
    ws.send(JSON.stringify({
      type: 'approval_response',
      token: pendingApprovalToken,
      approved: false,
      reason: reason || '用户已在界面上拒绝此操作'
    }));
    approvalCard.classList.add('hidden');
    pendingApprovalToken = null;
    updateStatus('thinking', '操作已拒绝，Agent 正在重新规划...');
  });

  // Antigravity Approval Mode Selector Handler
  if (approvalModeSelect) {
    approvalModeSelect.addEventListener('change', () => {
      const mode = approvalModeSelect.value;
      fetch('/api/approval-mode', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        body: JSON.stringify({ approval_mode: mode })
      }).catch(err => console.error('Failed to update approval mode', err));
    });
  }

  // -------------------------------------------------------------
  // Session Tree & Non-Destructive Time-Travel
  // -------------------------------------------------------------
  function populateSnapshots(snaps) {
    if (!sessionSnapshotSelect) return;
    sessionSnapshotSelect.innerHTML = `<option value="">🌳 分支树节点 (${snaps ? snaps.length : 0})</option>`;
    if (Array.isArray(snaps)) {
      snaps.forEach(s => {
        const opt = document.createElement('option');
        opt.value = s.turn_index;
        const branch = s.branch || 'main';
        const p = s.prompt || s.summary || '';
        const shortPrompt = p.length > 14 ? p.slice(0, 14) + '...' : p;
        const branchIcon = branch === 'main' ? '🌿' : '🔀';
        opt.textContent = `${branchIcon} [${branch}] #${s.turn_index + 1}: ${shortPrompt}`;
        sessionSnapshotSelect.appendChild(opt);
      });
    }
  }

  async function loadSnapshots() {
    if (!sessionSnapshotSelect) return;
    try {
      const res = await fetch('/api/session/snapshots');
      if (!res.ok) return;
      const snaps = await res.json();
      populateSnapshots(snaps);
    } catch (e) {
      console.error('loadSnapshots error', e);
    }
  }

  if (sessionSnapshotSelect) {
    sessionSnapshotSelect.addEventListener('change', async () => {
      const val = sessionSnapshotSelect.value;
      if (val === '') return;
      const turnIdx = parseInt(val, 10);
      if (confirm(`【会话分支树 • 时间旅行】确定要切换至该节点状态吗？\n所有已有分支与探索历史均被完整保存，您可以从此处随时发起新的分支生长，无任何历史被删除。`)) {
        try {
          const res = await fetch('/api/session/rollback', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json; charset=utf-8' },
            body: JSON.stringify({ turn_index: turnIdx })
          });
          const data = await res.json();
          if (res.ok) {
            // Smoothly redraw message tree from rollback response
            if (data.messages) {
              renderMessageTree(data.messages);
            }
            loadSnapshots();
            loadSessions();
          } else {
            alert('分支切换失败: ' + (data.error || '未知错误'));
          }
        } catch (e) {
          alert('分支切换异常: ' + e.message);
        }
      }
      sessionSnapshotSelect.value = '';
    });
  }

  // -------------------------------------------------------------
  // Git Status & Rollback Handlers
  // -------------------------------------------------------------
  async function checkGitStatus() {
    if (!gitWidget) return;
    try {
      const res = await fetch('/api/git/status');
      if (!res.ok) return;
      const data = await res.json();
      if (!data.is_git) {
        gitWidget.classList.add('hidden');
        gitWidget.classList.remove('flex');
        return;
      }
      gitWidget.classList.remove('hidden');
      gitWidget.classList.add('flex');
      const statusStr = (data.status || '').trim();
      if (statusStr) {
        const lines = statusStr.split('\n').filter(Boolean);
        gitStatusText.textContent = `Git: ${lines.length} 个改动`;
        gitStatusText.className = 'text-[#fdd663] font-mono text-[11px] font-medium';
        btnGitRollback.classList.remove('hidden');
      } else {
        gitStatusText.textContent = 'Git: 干净';
        gitStatusText.className = 'text-[#c4c7c5] font-mono text-[11px]';
        btnGitRollback.classList.add('hidden');
      }
    } catch (e) {
      console.error('checkGitStatus error', e);
    }
  }

  if (btnGitRollback) {
    btnGitRollback.addEventListener('click', async () => {
      if (!confirm('【安全回滚】确定要丢弃当前工作区所有未提交的代码变更吗？\n该操作将执行 git reset --hard 并清理所有未追踪的新建文件。')) return;
      try {
        const res = await fetch('/api/git/rollback', { method: 'POST' });
        const data = await res.json();
        if (res.ok) {
          checkGitStatus();
          alert('Git 工作区已秒级恢复至干净状态！\n' + (data.message || ''));
        } else {
          alert('Git 回滚失败: ' + (data.error || '未知错误'));
        }
      } catch (e) {
        alert('Git 回滚异常: ' + e.message);
      }
    });
  }

  // -------------------------------------------------------------
  // Workspace Switcher Handlers
  // -------------------------------------------------------------
  if (btnSwitchWorkspace && workspaceModal) {
    btnSwitchWorkspace.addEventListener('click', () => {
      wsNewPath.value = workspacePath.textContent.trim();
      wsSwitchError.classList.add('hidden');
      workspaceModal.classList.remove('hidden');
      wsNewPath.focus();
    });

    const closeWsModal = () => workspaceModal.classList.add('hidden');
    if (btnCloseWsModal) btnCloseWsModal.addEventListener('click', closeWsModal);
    if (btnCancelWs) btnCancelWs.addEventListener('click', closeWsModal);

    // Native Folder Browser Dialog Handler
    if (btnBrowseFolder) {
      btnBrowseFolder.addEventListener('click', async () => {
        btnBrowseFolder.disabled = true;
        btnBrowseFolderText.textContent = '正在唤起 Windows 文件资源管理器...';
        try {
          const res = await fetch('/api/workspace/browse', { method: 'POST' });
          const data = await res.json();
          if (data.status === 'ok' && data.path) {
            wsNewPath.value = data.path;
          }
        } catch (err) {
          console.error('Browse workspace error:', err);
        } finally {
          btnBrowseFolder.disabled = false;
          btnBrowseFolderText.textContent = '唤起系统文件资源管理器选择目录';
        }
      });
    }

    btnConfirmSwitchWs.addEventListener('click', () => {
      const newDir = wsNewPath.value.trim();
      if (!newDir) return;
      btnConfirmSwitchWs.disabled = true;
      btnConfirmSwitchWs.textContent = '切换中...';

      fetch('/api/workspace', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        body: JSON.stringify({ workspace: newDir })
      })
      .then(res => {
        if (!res.ok) {
          return res.text().then(t => { throw new Error(t); });
        }
        return res.json();
      })
      .then(data => {
        updateWorkspaceUI(data.workspace);
        closeWsModal();
        checkGitStatus();
      })
      .catch(err => {
        wsSwitchError.textContent = err.message || '切换失败';
        wsSwitchError.classList.remove('hidden');
      })
      .finally(() => {
        btnConfirmSwitchWs.disabled = false;
        btnConfirmSwitchWs.textContent = '切换并生效';
      });
    });
  }

  // Quick prompt cards
  document.querySelectorAll('.quick-prompt-card').forEach(card => {
    card.addEventListener('click', () => {
      const text = card.querySelector('p').textContent.trim();
      promptInput.value = text;
      promptInput.dispatchEvent(new Event('input'));
      promptInput.focus();
    });
  });

  // Auto-resize input textarea
  promptInput.addEventListener('input', () => {
    promptInput.style.height = 'auto';
    promptInput.style.height = Math.min(promptInput.scrollHeight, 140) + 'px';
  });

  function sendTask() {
    const prompt = promptInput.value.trim();
    if (!prompt) return;

    if (!ws || ws.readyState !== WebSocket.OPEN) {
      appendErrorMessage('⚠️ WebSocket 连接已断开或仍在建立中，正在为您重新连接，请稍候重试...');
      connectWS();
      return;
    }

    appendUserMessage(prompt);
    promptInput.value = '';
    promptInput.style.height = 'auto';

    activeMessageCard = null;
    currentBlockType = 'none';
    currentThoughtBox = null;
    currentThoughtContent = null;
    currentContentDiv = null;
    currentContentBuffer = '';

    const selectedModelId = modelSelect ? modelSelect.value : '';

    try {
      ws.send(JSON.stringify({
        type: 'user_prompt',
        prompt: prompt,
        model_provider: selectedModelId,
        thinking_level: currentThinkingLevel || 'high'
      }));
      updateStatus('thinking', '思考中...');
    } catch (err) {
      promptInput.value = prompt; // restore input on failure
      appendErrorMessage('指令发送失败: ' + err.message);
      updateStatus('idle', '发送失败');
    }
  }

  btnSend.addEventListener('click', sendTask);
  promptInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendTask();
    }
  });

  btnStop.addEventListener('click', () => {
    if (ws) {
      ws.send(JSON.stringify({ type: 'cancel' }));
      updateStatus('idle', '已取消任务');
    }
  });

  // Settings Modal Handlers
  if (btnSettings) {
    btnSettings.addEventListener('click', () => {
      settingsModal.classList.remove('hidden');
      loadModelProfiles();
      loadOfficialProviders();
    });
  }
  if (btnCloseSettings) btnCloseSettings.addEventListener('click', () => settingsModal.classList.add('hidden'));
  if (btnDoneSettings) btnDoneSettings.addEventListener('click', () => settingsModal.classList.add('hidden'));

  function escapeHtml(str) {
    if (!str) return '';
    return str.toString()
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }

  // Initialize
  loadConfig();
  loadOfficialProviders();
  connectWS();
});
