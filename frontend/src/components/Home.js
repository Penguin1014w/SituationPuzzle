import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function Home() {
  const navigate = useNavigate();
  const [scrollPosition, setScrollPosition] = useState(0);
  const [riddles, setRiddles] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/riddles')
      .then(response => response.json())
      .then(data => {
        setRiddles(data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching riddles:', error);
        setLoading(false);
      });
  }, []);

  const handleSelect = (id) => {
    navigate(`/game/${id}`);
  };

  const scrollLeft = () => {
    const container = document.querySelector('.riddle-list');
    const newPosition = scrollPosition - 300;
    setScrollPosition(Math.max(0, newPosition));
    container.scrollTo({
      left: Math.max(0, newPosition),
      behavior: 'smooth'
    });
  };

  const scrollRight = () => {
    const container = document.querySelector('.riddle-list');
    const newPosition = scrollPosition + 300;
    setScrollPosition(newPosition);
    container.scrollTo({
      left: newPosition,
      behavior: 'smooth'
    });
  };

  if (loading) {
    return (
      <div className="App">
        <div className="home-banner">
          <h1>情境谜题</h1>
          <p>加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <div className="home-banner">
        <h1>情境谜题</h1>
        <p>选择一个谜题开始游戏</p>
      </div>
      <div className="riddle-container">
        <button className="scroll-button left" onClick={scrollLeft}>&lt;</button>
        <div className="riddle-list">
          {riddles.map((riddle, index) => (
            <button
              key={riddle.id}
              className={`riddle-button ${index === scrollPosition / 300 ? 'selected' : ''}`}
              onClick={() => handleSelect(riddle.id)}
            >
              {riddle.title}
            </button>
          ))}
        </div>
        <button className="scroll-button right" onClick={scrollRight}>&gt;</button>
      </div>
    </div>
  );
}

export default Home;
