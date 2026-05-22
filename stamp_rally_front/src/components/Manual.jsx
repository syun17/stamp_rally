const Manual = ({ navigate, currentScreen }) => {
  return (
    <div className="screen">
      <div className="status-bar"><span>STAMP RALLY</span><span>●●●</span></div>
      <div className="nav-bar" style={{ background: '#00C060' }}>
        <div className="nav-back clickable" style={{ color: '#FFD900' }} onClick={() => navigate('home')}>‹ 戻る</div>
        <div className="nav-title" style={{ color: '#FFF' }}>遊び方</div>
        <div style={{ width: '36px' }}></div>
      </div>
      <div className="explain-body">
        <div className="step-card">
          <div className="step-num s1">1</div>
          <div>
            <div className="step-title">スタートする</div>
            <div className="step-desc">ホームの「スタート」を押してビンゴカードを受け取ろう！</div>
          </div>
        </div>
        <div className="step-card">
          <div className="step-num s2">2</div>
          <div>
            <div className="step-title">QRを読み取る</div>
            <div className="step-desc">各会場のQRをスキャンしてスタンプをゲット！</div>
          </div>
        </div>
        <div className="step-card">
          <div className="step-num s3">3</div>
          <div>
            <div className="step-title">ビンゴを目指す</div>
            <div className="step-desc">縦・横・斜めが揃ったらビンゴ達成！</div>
          </div>
        </div>
        <div className="step-card">
          <div className="step-num s4">4</div>
          <div>
            <div className="step-title">景品を受け取る</div>
            <div className="step-desc">景品画面をスタッフに見せよう！</div>
          </div>
        </div>
      </div>
      <div className="tab-bar">
        <div className={`tab-item clickable ${currentScreen === 'bingo' ? 'active' : ''}`} onClick={() => navigate('bingo')}>
          <svg className="tab-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <rect x="3" y="3" width="7" height="7" /><rect x="14" y="3" width="7" height="7" />
            <rect x="3" y="14" width="7" height="7" /><rect x="14" y="14" width="7" height="7" />
          </svg>
          ビンゴ
        </div>
        <div className={`tab-item clickable ${currentScreen === 'qr' ? 'active' : ''}`} onClick={() => navigate('qr')}>
          <svg className="tab-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <rect x="3" y="3" width="18" height="18" rx="2" /><path d="M9 9h.01M15 9h.01M9 15h.01M15 15h.01M12 12h.01" />
          </svg>
          QR
        </div>
        <div className={`tab-item clickable ${currentScreen === 'prize' ? 'active' : ''}`} onClick={() => navigate('prize')}>
          <svg className="tab-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <path d="M20 12V22H4V12" /><path d="M22 7H2v5h20V7z" /><path d="M12 22V7" />
          </svg>
          景品
        </div>
      </div>
    </div>
  );
};

export default Manual;