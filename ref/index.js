import React, { useState, useEffect } from 'react';
import { 
  ShoppingCart, 
  Search, 
  Menu, 
  X, 
  ChevronRight, 
  Minus, 
  Plus, 
  Trash2, 
  Phone, 
  Mail, 
  MapPin, 
  MessageCircle,
  ShieldCheck,
  TrendingUp,
  Award,
  Activity,
  TrendingDown
} from 'lucide-react';

// --- Mock Data ---
const PRODUCTS = [
  {
    id: 1,
    name: '加拿大楓葉白銀幣 1盎司',
    price: 980,
    category: '銀幣',
    image: 'https://images.unsplash.com/photo-1620055695026-664f3d1e6704?auto=format&fit=crop&q=80&w=400',
    description: '世界知名的投資銀幣，純度高達 99.99%。'
  },
  {
    id: 2,
    name: '澳洲袋鼠白銀幣 1盎司',
    price: 950,
    category: '銀幣',
    image: 'https://images.unsplash.com/photo-1596484552834-6a58f850d0fa?auto=format&fit=crop&q=80&w=400',
    description: '澳洲柏斯鑄幣廠出品，每年更新圖案，極具收藏價值。'
  },
  {
    id: 3,
    name: 'PAMP 財富女神白銀條塊 10盎司',
    price: 10500,
    category: '銀條',
    image: 'https://images.unsplash.com/photo-1610375461246-83df859d849d?auto=format&fit=crop&q=80&w=400',
    description: '瑞士頂級精煉廠製造，附有防偽認證與獨立序號。'
  },
  {
    id: 4,
    name: '瑞士純銀顆粒 1公斤',
    price: 32000,
    category: '工業用銀',
    image: 'https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?auto=format&fit=crop&q=80&w=400',
    description: '適合工業加工或大額避險投資，純度 99.9%。'
  }
];

export default function App() {
  const [activeTab, setActiveTab] = useState('home'); // 'home', 'cart', 'order'
  const [cart, setCart] = useState([]);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [toastMessage, setToastMessage] = useState('');

  // Toast Notification
  const showToast = (message) => {
    setToastMessage(message);
    setTimeout(() => setToastMessage(''), 3000);
  };

  // Cart Functions
  const addToCart = (product) => {
    setCart(prev => {
      const existing = prev.find(item => item.id === product.id);
      if (existing) {
        return prev.map(item => 
          item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        );
      }
      return [...prev, { ...product, quantity: 1 }];
    });
    showToast(`已將 ${product.name} 加入購物車`);
  };

  const updateQuantity = (id, delta) => {
    setCart(prev => prev.map(item => {
      if (item.id === id) {
        const newQuantity = Math.max(1, item.quantity + delta);
        return { ...item, quantity: newQuantity };
      }
      return item;
    }));
  };

  const removeFromCart = (id) => {
    setCart(prev => prev.filter(item => item.id !== id));
    showToast('已從購物車移除商品');
  };

  const cartTotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const cartItemCount = cart.reduce((sum, item) => sum + item.quantity, 0);

  // Navigation Helper
  const navigateTo = (tab) => {
    setActiveTab(tab);
    setIsMobileMenuOpen(false);
    window.scrollTo(0, 0);
  };

  return (
    <div className="min-h-screen flex flex-col font-sans bg-gray-50 text-gray-800">
      {/* Toast Notification */}
      {toastMessage && (
        <div className="fixed top-20 left-1/2 transform -translate-x-1/2 bg-slate-800 text-white px-6 py-3 rounded-lg shadow-xl z-50 transition-all duration-300">
          {toastMessage}
        </div>
      )}

      {/* Live Ticker */}
      <LiveTicker />

      {/* Navigation Bar */}
      <header className="bg-slate-900 text-white sticky top-0 z-40 shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            {/* Logo */}
            <div 
              className="flex items-center cursor-pointer" 
              onClick={() => navigateTo('home')}
            >
              <Award className="h-8 w-8 text-gray-300 mr-2" />
              <span className="text-2xl font-bold tracking-widest text-gray-100">SILVER LORD</span>
            </div>

            {/* Desktop Menu */}
            <nav className="hidden md:flex space-x-8 items-center">
              <button 
                onClick={() => navigateTo('home')}
                className={`hover:text-gray-300 transition-colors ${activeTab === 'home' ? 'text-white font-semibold border-b-2 border-white' : 'text-gray-400'}`}
              >
                首頁
              </button>
              <button 
                onClick={() => navigateTo('order')}
                className={`hover:text-gray-300 transition-colors ${activeTab === 'order' ? 'text-white font-semibold border-b-2 border-white' : 'text-gray-400'}`}
              >
                訂單查詢
              </button>
              <button 
                onClick={() => navigateTo('cart')}
                className="relative flex items-center hover:text-gray-300 transition-colors"
              >
                <ShoppingCart className="h-6 w-6" />
                {cartItemCount > 0 && (
                  <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                    {cartItemCount}
                  </span>
                )}
              </button>
            </nav>

            {/* Mobile Menu Button */}
            <div className="md:hidden flex items-center">
              <button 
                onClick={() => navigateTo('cart')}
                className="relative p-2 mr-2"
              >
                <ShoppingCart className="h-6 w-6" />
                {cartItemCount > 0 && (
                  <span className="absolute top-0 right-0 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                    {cartItemCount}
                  </span>
                )}
              </button>
              <button onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
                {isMobileMenuOpen ? <X className="h-7 w-7" /> : <Menu className="h-7 w-7" />}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Menu Dropdown */}
        {isMobileMenuOpen && (
          <div className="md:hidden bg-slate-800 px-4 pt-2 pb-4 space-y-2 shadow-inner">
            <button 
              onClick={() => navigateTo('home')}
              className="block w-full text-left px-3 py-3 rounded-md hover:bg-slate-700"
            >
              首頁
            </button>
            <button 
              onClick={() => navigateTo('order')}
              className="block w-full text-left px-3 py-3 rounded-md hover:bg-slate-700"
            >
              訂單查詢
            </button>
          </div>
        )}
      </header>

      {/* Main Content Area */}
      <main className="flex-grow">
        {activeTab === 'home' && <HomeView onAddToCart={addToCart} />}
        {activeTab === 'cart' && (
          <CartView 
            cart={cart} 
            updateQuantity={updateQuantity} 
            removeFromCart={removeFromCart} 
            cartTotal={cartTotal} 
            onCheckout={() => {
              if(cart.length === 0) return;
              showToast('導向結帳頁面中... (模擬功能)');
            }}
          />
        )}
        {activeTab === 'order' && <OrderQueryView />}
      </main>

      {/* Footer Section */}
      <footer className="bg-slate-900 text-gray-300 pt-12 pb-8 border-t border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            
            {/* 聯絡我們 */}
            <div>
              <h3 className="text-white text-lg font-bold mb-4 flex items-center">
                <Phone className="h-5 w-5 mr-2" /> 聯絡我們
              </h3>
              <ul className="space-y-3 text-sm">
                <li className="flex items-start">
                  <Phone className="h-4 w-4 mr-2 mt-1 flex-shrink-0" />
                  <span>客服專線：0927-066-022<br/><span className="text-gray-500">(週一至週五 09:00-18:00)</span></span>
                </li>
                <li className="flex items-center">
                  <Mail className="h-4 w-4 mr-2 flex-shrink-0" />
                  <span>service@silverlord.com.tw</span>
                </li>
                <li className="flex items-start">
                  <MapPin className="h-4 w-4 mr-2 mt-1 flex-shrink-0" />
                  <span>台北市信義區信義路五段7號 (台北101大樓)</span>
                </li>
              </ul>
            </div>

            {/* 匯款資訊 */}
            <div>
              <h3 className="text-white text-lg font-bold mb-4 flex items-center">
                <Award className="h-5 w-5 mr-2" /> 匯款資訊
              </h3>
              <div className="bg-slate-800 p-4 rounded-lg text-sm border border-slate-700">
                <p className="mb-2"><span className="text-gray-400">銀行名稱：</span>土地銀行 (005)</p>
                <p className="mb-2"><span className="text-gray-400">分行名稱：</span>中港分行</p>
                <p className="mb-2"><span className="text-gray-400">戶　　名：</span>吳仲凱</p>
                <p className="font-mono text-white text-lg tracking-wider mt-2 border-t border-slate-700 pt-2">
                  094005550186
                </p>
              </div>
            </div>

            {/* 訂購須知 */}
            <div>
              <h3 className="text-white text-lg font-bold mb-4 flex items-center">
                <ShieldCheck className="h-5 w-5 mr-2" /> 訂購須知
              </h3>
              <ul className="space-y-2 text-sm text-gray-400 list-disc list-inside">
                <li>貴金屬商品價格隨國際盤價波動，結帳金額以當下系統顯示為準。</li>
                <li>完成下單後請於 48 小時內完成匯款，逾期將自動取消訂單。</li>
                <li>依據消保法規定，實體黃金、白銀等貴金屬商品<span className="text-red-400 font-semibold">不適用七天鑑賞期</span>，售出後概不退換。</li>
                <li>現貨商品將於款項確認後 2 個工作天內物流寄出。</li>
              </ul>
            </div>

            {/* 官方LINE */}
            <div>
              <h3 className="text-white text-lg font-bold mb-4 flex items-center">
                <MessageCircle className="h-5 w-5 mr-2" /> 官方 LINE
              </h3>
              <p className="text-sm text-gray-400 mb-4">
                加入官方 LINE 帳號，獲得最新國際銀價資訊與專屬優惠通知，並提供一對一客服諮詢。
              </p>
              <button className="bg-[#06C755] hover:bg-[#05b34c] text-white px-4 py-2 rounded-lg font-bold flex items-center transition-colors w-full justify-center">
                <MessageCircle className="h-5 w-5 mr-2" />
                加入好友 @silverlord
              </button>
            </div>

          </div>
          
          <div className="border-t border-slate-800 mt-12 pt-8 text-center text-sm text-gray-500">
            &copy; {new Date().getFullYear()} SILVER LORD 貴金屬買賣平台. All rights reserved.
          </div>
        </div>
      </footer>
    </div>
  );
}

// --- Views ---

function LiveTicker() {
  const [quotes, setQuotes] = useState({
    silverUsd: 28.50,
    usdTwd: 32.15,
    trend: 'up'
  });

  // 模擬即時報價變動
  useEffect(() => {
    const interval = setInterval(() => {
      setQuotes(prev => {
        const changeUsd = (Math.random() - 0.45) * 0.05; // 隨機微幅波動
        const changeTwd = (Math.random() - 0.5) * 0.01;
        const newUsd = +(prev.silverUsd + changeUsd).toFixed(2);
        return {
          silverUsd: newUsd,
          usdTwd: +(prev.usdTwd + changeTwd).toFixed(2),
          trend: newUsd >= prev.silverUsd ? 'up' : 'down'
        };
      });
    }, 4000); // 每 4 秒更新一次
    return () => clearInterval(interval);
  }, []);

  const silverTwd = Math.round(quotes.silverUsd * quotes.usdTwd);

  return (
    <div className="bg-slate-950 text-gray-300 py-2 px-4 sm:px-6 lg:px-8 text-xs font-mono tracking-wide overflow-x-auto whitespace-nowrap">
      <div className="max-w-7xl mx-auto flex items-center justify-start md:justify-center space-x-6 min-w-max">
        <div className="flex items-center">
          <Activity className="h-3 w-3 mr-2 text-blue-400" />
          <span className="font-semibold text-white mr-2">即時市場報價</span>
        </div>
        <div className="flex items-center space-x-2">
          <span>XAG/USD (國際白銀) :</span>
          <span className={`font-bold flex items-center transition-colors duration-300 ${quotes.trend === 'up' ? 'text-green-400' : 'text-red-400'}`}>
            ${quotes.silverUsd.toFixed(2)}
            {quotes.trend === 'up' ? <TrendingUp className="h-3 w-3 ml-1" /> : <TrendingDown className="h-3 w-3 ml-1" />}
          </span>
        </div>
        <div className="text-gray-600">|</div>
        <div className="flex items-center space-x-2">
          <span>USD/TWD (美元/台幣) :</span>
          <span className="font-bold text-white transition-all duration-300">{quotes.usdTwd.toFixed(2)}</span>
        </div>
        <div className="text-gray-600">|</div>
        <div className="flex items-center space-x-2">
          <span>參考台幣銀價 :</span>
          <span className="font-bold text-yellow-400">NT$ {silverTwd} / oz</span>
        </div>
      </div>
    </div>
  );
}

function HomeView({ onAddToCart }) {
  return (
    <div className="animate-in fade-in duration-500">
      {/* Banner Area */}
      <section className="relative bg-slate-900 text-white overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src="https://images.unsplash.com/photo-1610375461246-83df859d849d?auto=format&fit=crop&q=80&w=1920" 
            alt="Silver Background" 
            className="w-full h-full object-cover opacity-30 mix-blend-overlay"
          />
          <div className="absolute inset-0 bg-gradient-to-r from-slate-900 via-slate-900/80 to-transparent"></div>
        </div>
        
        <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32">
          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-6">
            <span className="block text-gray-300 text-2xl md:text-3xl mb-2 font-normal">您的專屬實體白銀金庫</span>
            SILVER LORD
          </h1>
          <p className="text-lg md:text-xl max-w-xl mb-10 text-gray-300 leading-relaxed">
            我們提供國際認證、最高純度的白銀實體投資產品。透明的報價與安全的交易流程，是您避險保值與財富傳承的最佳選擇。
          </p>
          <div className="flex space-x-4">
            <button 
              onClick={() => {
                document.getElementById('products').scrollIntoView({ behavior: 'smooth' });
              }}
              className="bg-white text-slate-900 px-8 py-3 rounded-full font-bold hover:bg-gray-200 transition-colors shadow-lg"
            >
              選購白銀
            </button>
          </div>
        </div>
      </section>

      {/* 簡單介紹 Introduction */}
      <section className="py-16 bg-white border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-12">
            <h2 className="text-3xl font-bold text-slate-900 mb-4">為什麼選擇 SILVER LORD？</h2>
            <div className="w-16 h-1 bg-slate-900 mx-auto mb-6"></div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center p-6 rounded-2xl bg-gray-50 hover:shadow-md transition-shadow">
              <ShieldCheck className="h-12 w-12 mx-auto text-slate-700 mb-4" />
              <h3 className="text-xl font-bold mb-2">國際認證品質</h3>
              <p className="text-gray-600">所有銀條、銀幣皆來自 LBMA 認證之國際知名鑄幣廠，保證純度高達 99.9% 以上，附帶原廠防偽證明。</p>
            </div>
            <div className="text-center p-6 rounded-2xl bg-gray-50 hover:shadow-md transition-shadow">
              <TrendingUp className="h-12 w-12 mx-auto text-slate-700 mb-4" />
              <h3 className="text-xl font-bold mb-2">透明即時報價</h3>
              <p className="text-gray-600">價格緊貼國際白銀現貨盤價，公式透明公開，拒絕隱藏費用，讓您的投資效益最大化。</p>
            </div>
            <div className="text-center p-6 rounded-2xl bg-gray-50 hover:shadow-md transition-shadow">
              <Award className="h-12 w-12 mx-auto text-slate-700 mb-4" />
              <h3 className="text-xl font-bold mb-2">安心保值物流</h3>
              <p className="text-gray-600">採用高規格保值快遞配送，全程錄影監控並全額投保，確保您的實體資產安全送達手中。</p>
            </div>
          </div>
        </div>
      </section>

      {/* 商品方塊清單 Product Grid */}
      <section id="products" className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-end mb-10">
            <div>
              <h2 className="text-3xl font-bold text-slate-900 mb-2">熱銷白銀商品</h2>
              <p className="text-gray-500">實體白銀，掌握在手的真實財富</p>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {PRODUCTS.map(product => (
              <div key={product.id} className="bg-white rounded-2xl shadow-sm hover:shadow-xl transition-shadow overflow-hidden group flex flex-col">
                <div className="relative h-64 overflow-hidden bg-gray-100">
                  <img 
                    src={product.image} 
                    alt={product.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                  />
                  <div className="absolute top-4 left-4 bg-slate-900 text-white text-xs font-bold px-3 py-1 rounded-full">
                    {product.category}
                  </div>
                </div>
                <div className="p-6 flex-grow flex flex-col justify-between">
                  <div>
                    <h3 className="text-lg font-bold text-slate-900 mb-2 line-clamp-2">{product.name}</h3>
                    <p className="text-sm text-gray-500 mb-4 line-clamp-2">{product.description}</p>
                  </div>
                  <div>
                    <div className="text-2xl font-black text-slate-900 mb-4">
                      NT$ {product.price.toLocaleString()}
                    </div>
                    <button 
                      onClick={() => onAddToCart(product)}
                      className="w-full bg-slate-100 text-slate-900 hover:bg-slate-900 hover:text-white font-bold py-3 px-4 rounded-xl transition-colors flex justify-center items-center"
                    >
                      <ShoppingCart className="h-5 w-5 mr-2" />
                      加入購物車
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  );
}

function CartView({ cart, updateQuantity, removeFromCart, cartTotal, onCheckout }) {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 animate-in fade-in duration-300">
      <h2 className="text-3xl font-bold text-slate-900 mb-8 flex items-center">
        <ShoppingCart className="h-8 w-8 mr-3" />
        購物車清單
      </h2>

      {cart.length === 0 ? (
        <div className="bg-white rounded-2xl shadow-sm p-16 text-center border border-gray-100">
          <ShoppingCart className="h-16 w-16 mx-auto text-gray-300 mb-4" />
          <h3 className="text-xl font-bold text-gray-500 mb-2">您的購物車是空的</h3>
          <p className="text-gray-400 mb-6">快去商品區挑選具備投資價值的白銀吧！</p>
        </div>
      ) : (
        <div className="lg:grid lg:grid-cols-12 lg:gap-8">
          <div className="lg:col-span-8">
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
              <ul className="divide-y divide-gray-100">
                {cart.map((item) => (
                  <li key={item.id} className="p-6 flex flex-col sm:flex-row items-center sm:items-start gap-6">
                    <img 
                      src={item.image} 
                      alt={item.name} 
                      className="w-24 h-24 object-cover rounded-lg bg-gray-50"
                    />
                    <div className="flex-1 text-center sm:text-left w-full">
                      <div className="flex justify-between items-start">
                        <div>
                          <h4 className="text-lg font-bold text-slate-900">{item.name}</h4>
                          <p className="text-sm text-gray-500 mt-1">{item.category}</p>
                        </div>
                        <button 
                          onClick={() => removeFromCart(item.id)}
                          className="text-gray-400 hover:text-red-500 transition-colors p-2"
                          title="移除商品"
                        >
                          <Trash2 className="h-5 w-5" />
                        </button>
                      </div>
                      
                      <div className="mt-4 flex flex-col sm:flex-row justify-between items-center sm:items-end w-full">
                        <div className="flex items-center border border-gray-200 rounded-lg bg-gray-50 mb-4 sm:mb-0">
                          <button 
                            onClick={() => updateQuantity(item.id, -1)}
                            className="p-2 text-gray-600 hover:bg-gray-200 rounded-l-lg transition-colors"
                          >
                            <Minus className="h-4 w-4" />
                          </button>
                          <span className="w-12 text-center font-semibold text-slate-900">
                            {item.quantity}
                          </span>
                          <button 
                            onClick={() => updateQuantity(item.id, 1)}
                            className="p-2 text-gray-600 hover:bg-gray-200 rounded-r-lg transition-colors"
                          >
                            <Plus className="h-4 w-4" />
                          </button>
                        </div>
                        <div className="text-xl font-black text-slate-900">
                          NT$ {(item.price * item.quantity).toLocaleString()}
                        </div>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Checkout Summary */}
          <div className="lg:col-span-4 mt-8 lg:mt-0">
            <div className="bg-slate-900 text-white rounded-2xl shadow-lg p-6 sticky top-28">
              <h3 className="text-xl font-bold mb-6 border-b border-slate-700 pb-4">訂單摘要</h3>
              
              <div className="space-y-4 mb-6">
                <div className="flex justify-between text-gray-300">
                  <span>商品小計</span>
                  <span>NT$ {cartTotal.toLocaleString()}</span>
                </div>
                <div className="flex justify-between text-gray-300">
                  <span>保值物流運費</span>
                  <span>NT$ 150</span>
                </div>
              </div>
              
              <div className="border-t border-slate-700 pt-4 mb-8">
                <div className="flex justify-between items-end">
                  <span className="text-lg font-bold">總計金額</span>
                  <span className="text-3xl font-black text-white">
                    NT$ {(cartTotal + 150).toLocaleString()}
                  </span>
                </div>
                <p className="text-xs text-gray-400 mt-2 text-right">含稅價，運費已計入</p>
              </div>

              <button 
                onClick={onCheckout}
                className="w-full bg-white text-slate-900 font-bold py-4 px-4 rounded-xl hover:bg-gray-200 transition-colors flex justify-center items-center"
              >
                前往結帳 <ChevronRight className="h-5 w-5 ml-1" />
              </button>
              
              <div className="mt-4 text-xs text-gray-400 text-center flex items-center justify-center">
                <ShieldCheck className="h-4 w-4 mr-1" />
                安全加密結帳環境
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function OrderQueryView() {
  const [orderId, setOrderId] = useState('');
  const [phone, setPhone] = useState('');
  const [searchResult, setSearchResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSearch = (e) => {
    e.preventDefault();
    if (!orderId || !phone) return;
    
    setIsLoading(true);
    // 模擬 API 查詢延遲
    setTimeout(() => {
      setIsLoading(false);
      setSearchResult({
        id: orderId,
        date: new Date().toLocaleDateString(),
        status: '處理中 (等待匯款)',
        total: '32,150',
        shipping: '保值宅配'
      });
    }, 800);
  };

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16 animate-in fade-in duration-300 min-h-[60vh]">
      <div className="text-center mb-10">
        <h2 className="text-3xl font-bold text-slate-900 mb-4">訂單查詢</h2>
        <p className="text-gray-500">請輸入您的訂單編號與訂購人手機號碼以查詢最新進度。</p>
      </div>

      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 md:p-8">
        <form onSubmit={handleSearch} className="space-y-6">
          <div>
            <label htmlFor="orderId" className="block text-sm font-medium text-gray-700 mb-2">
              訂單編號
            </label>
            <input
              type="text"
              id="orderId"
              value={orderId}
              onChange={(e) => setOrderId(e.target.value)}
              placeholder="例如：SL20260408XXXX"
              className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-slate-900 focus:border-slate-900 outline-none transition-all"
              required
            />
          </div>
          <div>
            <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-2">
              訂購人手機號碼
            </label>
            <input
              type="tel"
              id="phone"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="例如：0912345678"
              className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-slate-900 focus:border-slate-900 outline-none transition-all"
              required
            />
          </div>
          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-slate-900 text-white font-bold py-3 px-4 rounded-lg hover:bg-slate-800 transition-colors flex justify-center items-center disabled:bg-slate-400"
          >
            {isLoading ? (
              <span className="animate-pulse">查詢中...</span>
            ) : (
              <>
                <Search className="h-5 w-5 mr-2" />
                查詢訂單
              </>
            )}
          </button>
        </form>

        {/* 查詢結果 */}
        {searchResult && (
          <div className="mt-10 pt-8 border-t border-gray-100 animate-in slide-in-from-bottom-4 duration-500">
            <h3 className="text-lg font-bold text-slate-900 mb-4">查詢結果</h3>
            <div className="bg-gray-50 rounded-lg p-6 border border-gray-200">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-500">訂單編號</p>
                  <p className="font-mono font-medium text-slate-900">{searchResult.id}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">訂購日期</p>
                  <p className="font-medium text-slate-900">{searchResult.date}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">訂單狀態</p>
                  <p className="font-medium text-orange-600 font-bold">{searchResult.status}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">訂單總金額</p>
                  <p className="font-medium text-slate-900">NT$ {searchResult.total}</p>
                </div>
              </div>
              
              <div className="mt-6 p-4 bg-orange-50 border border-orange-100 rounded-md">
                <p className="text-sm text-orange-800">
                  <strong>溫馨提醒：</strong> 您的訂單正在等待匯款，請於下單後 24 小時內完成匯款，以免訂單被系統自動取消。匯款資訊請參考網頁下方說明。
                </p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}