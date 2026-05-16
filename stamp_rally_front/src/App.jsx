import { useState, useEffect, useRef } from 'react'
import './App.css'

const USER_ID = 'user_001'

// ビンゴグリッドの順番 (3x3, spot_id 1〜9)
const GRID = [
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9],
]

// ビンゴライン定義
const BINGO_LINES = [
  [1, 2, 3], [4, 5, 6], [7, 8, 9],
  [1, 4, 7], [2, 5, 8], [3, 6, 9],
  [1, 5, 9], [3, 5, 7],
]

function isBingoLine(line, stampedSet) {
  return line.every((id) => stampedSet.has(id))
}

function App() {
  const [spots, setSpots] = useState([])
  const [bingo, setBingo] = useState({ stamped_ids: [], bingo_count: 0, bingo_lines: [], is_complete: false })
  const [message, setMessage] = useState('')
  const [nfcSupported, setNfcSupported] = useState(false)
  const [nfcScanning, setNfcScanning] = useState(false)
  const nfcAbortRef = useRef(null)

  // スポット一覧を取得
  useEffect(() => {
    fetch('/api/spots')
      .then((res) => res.json())
      .then((data) => setSpots(data))
      .catch((err) => setMessage('スポット取得失敗: ' + err.message))
  }, [])

  // ビンゴ状況を取得
  const fetchBingo = () => {
    fetch(`/api/bingo/${USER_ID}`)
      .then((res) => res.json())
      .then((data) => setBingo(data))
      .catch((err) => setMessage('ビンゴ状況取得失敗: ' + err.message))
  }

  useEffect(() => {
    fetchBingo()
    setNfcSupported('NDEFReader' in window)
  }, [])

  // スタンプ取得（共通）
  const acquireStamp = async (body) => {
    const res = await fetch('/api/stamps/nfc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const data = await res.json()
    if (res.status === 201) {
      setMessage('スタンプを取得しました！')
      fetchBingo()
    } else if (res.status === 409) {
      setMessage('このスポットはすでにスタンプ済みです')
    } else {
      setMessage('エラー: ' + (data.error ?? '不明なエラー'))
    }
  }

  // NFCスキャン開始
  const startNfcScan = async () => {
    if (!nfcSupported) return
    try {
      const ndef = new window.NDEFReader()
      const controller = new AbortController()
      nfcAbortRef.current = controller

      await ndef.scan({ signal: controller.signal })
      setNfcScanning(true)
      setMessage('NFCタグをスマートフォンにかざしてください')

      ndef.addEventListener('reading', ({ serialNumber }) => {
        // シリアルナンバー（UID）をコロン区切り大文字に正規化
        const uid = serialNumber.toUpperCase().replace(/-/g, ':')
        setMessage(`NFCタグ検出: ${uid}`)
        acquireStamp({ user_id: USER_ID, nfc_uid: uid })
      })
    } catch (err) {
      setMessage('NFCスキャン失敗: ' + err.message)
      setNfcScanning(false)
    }
  }

  // NFCスキャン停止
  const stopNfcScan = () => {
    if (nfcAbortRef.current) {
      nfcAbortRef.current.abort()
      nfcAbortRef.current = null
    }
    setNfcScanning(false)
    setMessage('NFCスキャンを停止しました')
  }

  const stampedSet = new Set(bingo.stamped_ids ?? [])
  const completedLineSet = new Set((bingo.bingo_lines ?? []).flat())

  const spotMap = Object.fromEntries(spots.map((s) => [s.id, s]))

  return (
    <div style={{ maxWidth: 480, margin: '2rem auto', fontFamily: 'sans-serif', padding: '0 1rem' }}>
      <h1 style={{ textAlign: 'center' }}>スタンプラリー ビンゴ</h1>
      <p style={{ textAlign: 'center', color: '#555' }}>ユーザーID: <strong>{USER_ID}</strong></p>

      {/* ビンゴ達成表示 */}
      {bingo.bingo_count > 0 && (
        <div style={{ background: '#fff3cd', border: '2px solid #ffc107', borderRadius: 8, padding: '0.75rem 1rem', marginBottom: '1rem', textAlign: 'center', fontWeight: 'bold' }}>
          {bingo.is_complete ? '🎉 全スポット制覇！' : `🎯 ビンゴ ${bingo.bingo_count} ライン達成！`}
        </div>
      )}

      {/* メッセージ */}
      {message && (
        <div style={{ background: '#f0f0f0', padding: '0.5rem 1rem', borderRadius: 4, marginBottom: '1rem', fontSize: '0.9rem' }}>
          {message}
        </div>
      )}

      {/* NFCボタン */}
      {nfcSupported ? (
        <div style={{ textAlign: 'center', marginBottom: '1.5rem' }}>
          {nfcScanning ? (
            <button onClick={stopNfcScan} style={{ background: '#dc3545', color: '#fff', border: 'none', borderRadius: 8, padding: '0.75rem 2rem', fontSize: '1rem', cursor: 'pointer' }}>
              NFCスキャン停止
            </button>
          ) : (
            <button onClick={startNfcScan} style={{ background: '#0d6efd', color: '#fff', border: 'none', borderRadius: 8, padding: '0.75rem 2rem', fontSize: '1rem', cursor: 'pointer' }}>
              NFCスキャン開始
            </button>
          )}
        </div>
      ) : (
        <p style={{ textAlign: 'center', color: '#999', fontSize: '0.85rem' }}>
          ※ このブラウザはWeb NFCに未対応です（Android Chrome が必要）
        </p>
      )}

      {/* ビンゴグリッド */}
      {spots.length > 0 && (
        <div>
          <h2 style={{ textAlign: 'center', marginBottom: '0.5rem' }}>ビンゴカード</h2>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8 }}>
            {GRID.flat().map((spotId) => {
              const spot = spotMap[spotId]
              const stamped = stampedSet.has(spotId)
              const onBingoLine = completedLineSet.has(spotId)
              return (
                <div
                  key={spotId}
                  style={{
                    border: onBingoLine ? '2px solid #ffc107' : '1px solid #ccc',
                    borderRadius: 8,
                    padding: '0.75rem 0.5rem',
                    textAlign: 'center',
                    background: stamped ? (onBingoLine ? '#fff3cd' : '#e6f4ea') : '#fafafa',
                    position: 'relative',
                  }}
                >
                  <div style={{ fontSize: '0.75rem', color: '#666', marginBottom: 4 }}>
                    No.{spotId}
                  </div>
                  <div style={{ fontSize: '0.85rem', fontWeight: 'bold', lineHeight: 1.3 }}>
                    {spot ? spot.name : '?'}
                  </div>
                  {stamped && (
                    <div style={{ fontSize: '1.5rem', position: 'absolute', top: 4, right: 6 }}>
                      ✅
                    </div>
                  )}
                </div>
              )
            })}
          </div>
          <p style={{ textAlign: 'center', marginTop: '0.75rem', color: '#555' }}>
            取得済み: {stampedSet.size} / {spots.length} スポット
          </p>
        </div>
      )}
    </div>
  )
}

export default App

