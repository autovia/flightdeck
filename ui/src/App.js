// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import { useState, useEffect } from 'react'
import { Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Namespace from './components/Namespace';

function App() {

   const [isLoggedIn, setIsLoggedIn] = useState(false);
   const [isLoaded, setIsLoaded] = useState(false);
   
   useEffect(() => {
      fetch('/api/v1/auth/status', {
         method: 'GET',
      })
     .then(res => res)
      .then(d => {
         if (d.status == 200) {
            setIsLoggedIn(true);
            console.log('/api/v1/auth/status', d.status);
         }
         setIsLoaded(true);
      });
   }, []);

   return (
      <>{ isLoaded ? 
         <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={!isLoggedIn ? <Login /> : <Namespace />} />
            <Route path="/cluster/:cluster" element={!isLoggedIn ? <Login /> : <Namespace />} />
            <Route path="/node/:node" element={!isLoggedIn ? <Login /> : <Namespace />} />
            <Route path="/namespace/:namespace/:kind" element={!isLoggedIn ? <Login /> : <Namespace />} />
            <Route path="/resource/:namespace/:kind/:resource" element={!isLoggedIn ? <Login /> : <Namespace />} />  
         </Routes>
         : null}
      </>   
   );
};

export default App;