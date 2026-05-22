import '../App.css';


const Qr = ({ navigate, currentScreen }) => {
  return (
    <div className="screen">
      <div className="status-bar"><span>STAMP RALLY</span><span>●●●</span></div>
      <div className="nav-bar" style={{ background: '#00D4D4' }}>
        <div className="nav-back clickable" onClick={() => navigate('bingo')}>‹ 戻る</div>
        <div className="nav-title">QR読み込み</div>
        <div style={{ width: '36px' }}></div>
      </div>
      <div style={{ background: '#FAFAFA', flex: 1, display: 'flex', flexDirection: 'column' }}>
        <div className="qr-area">
          <div className="qr-frame">
            <div className="qr-mock">
              <div className="qr-px"></div><div className="qr-px"></div><div></div><div className="qr-px"></div><div className="qr-px"></div>
              <div className="qr-px"></div><div></div><div className="qr-px"></div><div></div><div className="qr-px"></div>
              <div></div><div className="qr-px"></div><div className="qr-px"></div><div className="qr-px"></div><div></div>
              <div className="qr-px"></div><div></div><div className="qr-px"></div><div></div><div className="qr-px"></div>
              <div className="qr-px"></div><div className="qr-px"></div><div></div><div className="qr-px"></div><div className="qr-px"></div>
            </div>
            <div className="scan-line-mock"></div>
          </div>
          <div className="qr-label">📷 枠内にQRを合わせてね！</div>
        </div>
        <div className="info-strip">
          <svg className="info-strip-icon" viewBox="0 0 24 24" fill="none" stroke="#111" strokeWidth="2.5" strokeLinecap="round">
            <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <div className="info-strip-text">各会場のQRコードをスキャンするとスタンプがもらえます！</div>
        </div>
        <div style={{ flex: 1 }}></div>
        <div className="big-btn outline clickable" onClick={() => navigate('bingo')}>← 戻る</div>
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

export default Qr;