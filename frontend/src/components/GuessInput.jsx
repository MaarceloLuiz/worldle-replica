import { useState, useEffect, useRef } from 'react';

const GuessInput = ({ territories, onSubmit, disabled }) => {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const inputRef = useRef(null);
  const suggestionsRef = useRef(null);

  const getSuggestions = (value) => {
    const query = value.toLowerCase();
    // If territories is not an array or is empty, return empty array
    if (!Array.isArray(territories) || territories.length === 0) {
      console.log("Territories not available:", territories);
      return [];
    }

    // Check if territories are strings or objects
    const isStringArray = typeof territories[0] === 'string';
    
    // If input is empty, return first 8 territories
    if (!query) {
      return territories.slice(0, 8);
    }
    
    // Filter based on whether we have strings or objects
    if (isStringArray) {
      return territories
        .filter(t => t && t.toLowerCase().includes(query))
        .slice(0, 8);
    } else {
      return territories
        .filter(t => t && t.name)
        .filter(t =>
          t.name.toLowerCase().includes(query) ||
          (t.code && t.code.toLowerCase() === query)
        )
        .slice(0, 8);
    }
  };

  const handleSubmit = (territory) => {
    if (!territory || disabled) return;
    
    // Check if territory is a string or object
    const territoryValue = typeof territory === 'string' ? territory : territory.name;
    
    onSubmit(territoryValue);
    setInput('');
    setSuggestions([]);
    setShowSuggestions(false);
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      try {
        const newSuggestions = getSuggestions(input);
        setSuggestions(newSuggestions);
        setShowSuggestions(newSuggestions.length > 0);
      } catch (error) {
        console.error("Error getting suggestions:", error);
        setSuggestions([]);
        setShowSuggestions(false);
      }
    }, 200);
    return () => clearTimeout(timer);
  }, [input, territories]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      if (selectedIndex >= 0 && selectedIndex < suggestions.length) {
        // Use the selected suggestion
        const selected = suggestions[selectedIndex];
        handleSubmit(selected);
      } else if (suggestions.length > 0) {
        // If no selection but we have suggestions, use the first one
        handleSubmit(suggestions[0]);
      } else if (input.trim()) {
        // If no suggestions but we have input, try to submit it
        handleSubmit(input);
      }
    }
    if (e.key === 'ArrowDown') {
      setSelectedIndex(prev => Math.min(prev + 1, suggestions.length - 1));
    }
    if (e.key === 'ArrowUp') {
      setSelectedIndex(prev => Math.max(prev - 1, -1));
    }
  };

  // Handle clicks outside the component to close suggestions
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        inputRef.current && 
        !inputRef.current.contains(event.target) && 
        suggestionsRef.current && 
        !suggestionsRef.current.contains(event.target)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="guess-input-container">
      <input
        ref={inputRef}
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        onFocus={() => {
          // Show all territories when focused with empty input
          const newSuggestions = getSuggestions(input);
          setSuggestions(newSuggestions);
          setShowSuggestions(true); // Always show on focus
        }}
        placeholder="Start typing country name..."
        aria-label="Country guess input"
        disabled={disabled}
        autoFocus
      />
      {showSuggestions && suggestions.length > 0 && (
        <ul ref={suggestionsRef} className="suggestions-list">
          {suggestions.map((territory, index) => {
            // Handle both string and object territories
            const displayText = typeof territory === 'string' ? territory : territory.name;
            
            return (
              <li
                key={index}
                className={`suggestion-item ${index === selectedIndex ? 'selected' : ''}`}
                onClick={() => handleSubmit(territory)}
              >
                {displayText}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
};

export default GuessInput;