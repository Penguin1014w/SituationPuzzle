// 导入必要的React组件和hooks
import React, { useState, useEffect } from 'react';
import './App.css';

// 主应用组件
function App() {
  // 状态管理
  const [riddles, setRiddles] = useState([]); // 存储所有谜面数据
  const [selectedRiddle, setSelectedRiddle] = useState(null); // 当前选中的谜面
  const [guess, setGuess] = useState(''); // 用户输入的猜测
  const [attempts, setAttempts] = useState(0); // 猜测次数
  const [gameOver, setGameOver] = useState(false); // 游戏是否结束
  const [loading, setLoading] = useState(true); // 加载状态
  const [chatHistory, setChatHistory] = useState([]); // 聊天历史记录
  const [scrollPosition, setScrollPosition] = useState(0); // 谜面列表滚动位置
  const [lang, setLang] = useState('CH'); // 当前语言，默认中文
  const riddleCardWidth = 320; // 谜面卡片宽度（包含间距）
  const [showAnswer, setShowAnswer] = useState(false); // 添加显示谜底的状态

  // 组件加载时获取谜面数据
  useEffect(() => {
    if (selectedRiddle !== null) return; // 只有在未选中谜面时才请求
    const fetchRiddles = async () => {
      try {
        setLoading(true);
        const response = await fetch(`https://situationpuzzle.onrender.com/api/riddles?lang=${lang}`);
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        setRiddles(data);
      } catch (error) {
        console.error('获取谜面失败，请检查后端服务是否运行', error);
      } finally {
        setLoading(false);
      }
    };
    fetchRiddles();
  }, [lang, selectedRiddle]);

  // 处理谜面选择
  const handleRiddleSelect = (riddle) => {
    // 确保保存完整的谜面数据，包括谜底
    const selectedRiddleData = {
      ...riddle,
      answer_ch: riddle.answer_ch,
      answer_en: riddle.answer_en,
      difficulty: riddle.difficulty
    };
    setSelectedRiddle(selectedRiddleData);
    setGuess('');
    setAttempts(0);
    setGameOver(false);
    setChatHistory([]);
    setShowAnswer(false); // 重置显示谜底状态
  };

  // 处理用户提交的猜测
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (attempts >= 5) {
      setGameOver(true);
      console.log(lang === 'CH' ? '游戏结束！你已经用完所有猜测次数。' : 'Game Over! You have used all your attempts.');
      return;
    }
    setChatHistory(prev => [...prev, { role: 'user', content: guess }]);
    try {
      // 发送猜测到后端进行验证
      const response = await fetch('https://situationpuzzle.onrender.com/api/check-answer', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          riddleId: selectedRiddle.id,
          answer: guess,
          lang: lang,
        }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setAttempts(prev => prev + 1);
      let aiReply = '';
      // 根据后端返回的状态生成相应的回复
      switch (data.status) {
        case 'correct':
          setGameOver(true);
          aiReply = lang === 'CH' ? '恭喜你猜对了！' : 'Congratulations! You got it right!';
          break;
        case 'yes':
          aiReply = lang === 'CH' ? '是的！' : 'Yes!';
          break;
        case 'no':
          aiReply = lang === 'CH' ? '不是！' : 'No!';
          break;
        case 'irrelevant':
          aiReply = lang === 'CH' 
            ? `你的问题与题目无关，请仔细思考谜面中的关键信息。还剩${5 - attempts - 1}次机会。`
            : `Your question is irrelevant to the puzzle. Please think about the key information in the riddle. ${5 - attempts - 1} attempts remaining.`;
          break;
        default:
          aiReply = lang === 'CH' ? '发生错误：未知的答案状态' : 'Error: Unknown answer status';
      }
      setChatHistory(prev => [...prev, { role: 'ai', content: aiReply }]);
      setGuess('');
    } catch (error) {
      console.error('发生错误：', error.message);
      setChatHistory(prev => [...prev, { role: 'ai', content: `发生错误：${error.message}` }]);
    }
  };

  // 处理谜面列表向左滚动
  const scrollLeft = () => {
    const container = document.querySelector('.riddle-list');
    const newPosition = Math.max(0, scrollPosition - riddleCardWidth);
    setScrollPosition(newPosition);
    container.scrollTo({ left: newPosition, behavior: 'smooth' });
  };

  // 处理谜面列表向右滚动
  const scrollRight = () => {
    const container = document.querySelector('.riddle-list');
    const maxScroll = (riddles.length - 1) * riddleCardWidth;
    const newPosition = Math.min(scrollPosition + riddleCardWidth, maxScroll);
    setScrollPosition(newPosition);
    container.scrollTo({ left: newPosition, behavior: 'smooth' });
  };

  // 切换语言
  const toggleLang = () => {
    setLang(lang === 'CH' ? 'EN' : 'CH');
  };

  // 渲染主界面
  return (
    <div className="App">
      <header className="App-header">
        <h1>Turtle Yes or No海龟汤</h1>
        <button className="lang-menu-btn" onClick={toggleLang} title="切换语言">
          {lang}
        </button>
      </header>
      <main>
        {lang === 'CH' ? (
          <div className="rules-box">
            <h2>游戏规则</h2>
            <p>
              海龟汤是一种情境推理游戏。主持人给出一个离奇的情境，玩家通过提问来还原真相。<br />
              你可以提出任何问题，主持人只会回答"是"、"不是"或"无关"。<br />
              目标是通过有限的提问，推理出完整的故事真相！你会有5次机会提问。<br />
            </p>
          </div>
        ) : (
          <div className="rules-box">
            <h2>Rules</h2>
            <p>
              Situation Puzzle (Lateral Thinking Puzzle) is a reasoning game. The host gives a bizarre scenario, and players ask questions to uncover the truth.<br />
              You can ask any question, and the host will only answer "Yes", "No", or "Irrelevant".<br />
              The goal is to deduce the full story with limited questions! you have 5 attempts.
            </p>
          </div>
        )}
        {loading ? (
          <p>{lang === 'CH' ? '加载中...' : 'Loading...'}</p>
        ) : !selectedRiddle ? (
          // 谜面选择界面
          <div className="riddle-container">
            <button className="scroll-button left" onClick={scrollLeft}>&lt;</button>
            <div className="riddle-list" style={{overflowX: 'hidden', display: 'flex', alignItems: 'center', gap: '30px', padding: '40px 0', flex: 1}}>
              {riddles.length > 0 ? (
                riddles.map((riddle, idx) => (
                  <div
                    key={riddle.id}
                    className={`riddle-button${Math.round(scrollPosition / riddleCardWidth) === idx ? ' selected' : ''}`}
                    style={{minWidth: '300px', maxWidth: '320px', height: '200px', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', cursor: 'pointer'}}
                    onClick={() => handleRiddleSelect(riddle)}
                  >
                    <h3 style={{marginBottom: 10}}>{riddle.title}</h3>
                    <p style={{fontSize: '1em', color: '#222'}}>{riddle.content}</p>
                    <div className="difficulty">
                      {Array(Math.max(1, riddle.difficulty)).fill('⭐').join('')}
                    </div>
                  </div>
                ))
              ) : (
                <p>{lang === 'CH' ? '没有可用的谜面' : 'No riddles available'}</p>
              )}
            </div>
            <button className="scroll-button right" onClick={scrollRight}>&gt;</button>
          </div>
        ) : (
          // 游戏界面
          <div className="game-container">
            <button className="back-btn" onClick={() => setSelectedRiddle(null)} style={{ marginBottom: 16 }}>
              {lang === 'CH' ? '返回' : 'Back'}
            </button>
            <h2>{selectedRiddle.title}</h2>
            <p>{selectedRiddle.content}</p>
            <div className="difficulty">
              {Array(Math.max(1, selectedRiddle.difficulty)).fill('⭐').join('')}
            </div>
            <div className="draft-chat-history">
              {chatHistory.map((msg, idx) => (
                <div
                  key={idx}
                  className={`draft-box draft-msg ${msg.role === 'user' ? 'draft-user' : 'draft-ai'}`}
                  style={{ textAlign: msg.role === 'user' ? 'right' : 'left' }}
                >
                  {msg.content}
                </div>
              ))}
            </div>
            {!gameOver ? (
              <form onSubmit={handleSubmit} className="draft-form">
                <input
                  id="guess-input"
                  name="guess"
                  type="text"
                  value={guess}
                  onChange={(e) => setGuess(e.target.value)}
                  placeholder={lang === 'CH' ? "输入你的猜测或提问..." : "Enter your guess or question..."}
                  className="draft-input"
                />
                <button type="submit" className="draft-btn">{lang === 'CH' ? '提交' : 'Submit'}</button>
              </form>
            ) : (
              <div className="game-over-buttons">
                {!showAnswer ? (
                  <button onClick={() => setShowAnswer(true)} className="draft-btn show-answer-btn">
                    {lang === 'CH' ? '查看谜底' : 'Show Answer'}
                  </button>
                ) : (
                  <div className="answer-box">
                    <h3>{lang === 'CH' ? '谜底：' : 'Answer:'}</h3>
                    <p>{lang === 'CH' ? selectedRiddle.answer_ch : selectedRiddle.answer_en}</p>
                  </div>
                )}
                <button onClick={() => setSelectedRiddle(null)} className="draft-btn back-btn">
                  {lang === 'CH' ? '返回谜面列表' : 'Back to Riddles'}
                </button>
              </div>
            )}
          </div>
        )}
      </main>
      <footer className="footer">
        <div>
          © 2025 Penguin1014w. All rights reserved.
          <br />
          Made by <a href="https://github.com/Penguin1014w/SituationPuzzle" target="_blank" rel="noopener noreferrer">Penguin1014w</a>
        </div>
      </footer>
    </div>
  );
}

export default App;