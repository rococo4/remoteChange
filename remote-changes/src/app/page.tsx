"use client";

import { useState, useEffect } from "react";

export default function Page() {
  const [activeTab, setActiveTab] = useState("Team");
  const [selectedConfig, setSelectedConfig] = useState<number | null>(null);
  const [editContent, setEditContent] = useState("");
  // NEW: State to track editing mode
  const [isEditing, setIsEditing] = useState(false);

  // New states for admin team management
  const [adminTeam, setAdminTeam] = useState<{ name: string; users: string[] } | null>(null);
  const [newTeamName, setNewTeamName] = useState("");
  const [newUserName, setNewUserName] = useState("");

  // NEW: States for Create functionality on Team tab
  const [showCreate, setShowCreate] = useState(false);
  const [newConfig, setNewConfig] = useState("");
  // NEW: state for description
  const [newConfigDescription, setNewConfigDescription] = useState("");

  // NEW: state для роли пользователя
  const [role, setRole] = useState("");

  // Add new state for authentication mode
  const [authMode, setAuthMode] = useState("login");

  // New states for authentication
  const [loginUsername, setLoginUsername] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [registerUsername, setRegisterUsername] = useState("");
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");

  // NEW: State for JWT token
  const [jwtToken, setJwtToken] = useState<string | null>(null);

  // Add new state for team info and team configs
  const [teamInfo, setTeamInfo] = useState<{ id: number; name: string } | null>(null);
  const [teamConfigs, setTeamConfigs] = useState<Array<{ id: number; name: string; type: string; date: string; description: string }>>([]);

  // NEW: Add state for config changes history
  const [configHistory, setConfigHistory] = useState<Array<{ username: string; action: string; date: string }>>([]);

  // NEW state for team users
  const [teamUsers, setTeamUsers] = useState<any[]>([]);
  const [newUserToAdd, setNewUserToAdd] = useState("");

  // NEW: Update useEffect to include token in headers and check response
  useEffect(() => {
    const storedToken = localStorage.getItem("jwtToken");
    setJwtToken(storedToken);
    fetch("http://localhost:8080/teams/user/role", {
      headers: { "Authorization": "Bearer " + storedToken }
    })
      .then((res) => {
        if (!res.ok) return res.text().then(text => { throw new Error(text) });
        return res.json();
      })
      .then((data) => {
        // data is a string: "admin" or "user"
        setRole(data);
        if (data === "admin") {
          setActiveTab("Admin");
        }
      })
      .catch((err) => console.error("Error fetching role:", err));
  }, []);

  // Add useEffect to fetch team info and configs when Team tab is active
  useEffect(() => {
    if (activeTab === "Team" && jwtToken) {
      fetch("http://localhost:8080/teams/user/team", {
        headers: { "Authorization": "Bearer " + jwtToken }
      })
        .then((res) => {
          if (!res.ok) return res.text().then(text => { throw new Error("Failed to fetch team info: " + text) });
          return res.json();
        })
        .then((data) => {
          setTeamInfo(data);
          return fetch(`http://localhost:8080/configs/team/${data.id}`, {
            headers: { "Authorization": "Bearer " + jwtToken }
          });
        })
        .then((res) => {
          if (!res.ok) return res.text().then(text => { throw new Error("Failed to fetch team configs: " + text) });
          return res.json();
        })
        .then((configs) => setTeamConfigs(configs))
        .catch((err) => console.error(err));
    }
  }, [activeTab, jwtToken]);

  // Add a new useEffect to fetch users when teamInfo exists and user is admin
  useEffect(() => {
    if (teamInfo && jwtToken) {
      fetch(`http://localhost:8080/teams/${teamInfo.id}/users`, {
        headers: { "Authorization": "Bearer " + jwtToken }
      })
        .then((res) => {
          if (!res.ok) throw new Error("Failed to fetch users");
          return res.json();
        })
        .then((users: string[]) => setTeamUsers(users))
        .catch((err) => console.error(err));
    }
  }, [teamInfo, jwtToken]);

  // Handler functions for login and registration
  const handleLogin = async () => {
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: loginUsername, password: loginPassword }),
      });
      const data = await res.json();
      if (data.token) {
        localStorage.setItem("jwtToken", data.token);
        setJwtToken(data.token);
        console.log("Login successful");
      } else {
        console.error("Login failed", data);
      }
    } catch (error) {
      console.error("Error during login", error);
    }
  };

  const handleRegister = async () => {
    try {
      const res = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: registerUsername,
          email: registerEmail,
          password: registerPassword,
        }),
      });
      const data = await res.json();
      if (data.token) {
        localStorage.setItem("jwtToken", data.token);
        setJwtToken(data.token);
        console.log("Registration successful");
      } else {
        console.error("Registration failed", data);
      }
    } catch (error) {
      console.error("Error during registration", error);
    }
  };

  return (
    <div className="flex h-screen font-sans">
      {/* Sidebar */}
      <div className="w-64 bg-gradient-to-b from-blue-600 to-blue-800 text-white p-6 shadow-lg">
        <ul className="space-y-6">
          {["Team", "Admin", "Info", "Login/Register"].map((tab) => (
            // Скрываем вкладку Admin, если роль не admin
            tab === "Admin" && role !== "admin" ? null : (
              <li
                key={tab}
                onClick={() => {
                  setActiveTab(tab);
                  setSelectedConfig(null);
                }}
                className={`cursor-pointer px-4 py-2 rounded-md transition-colors hover:bg-blue-500 ${
                  activeTab === tab && "bg-blue-400"
                }`}
              >
                {tab}
              </li>
            )
          ))}
        </ul>
      </div>

      {/* Main Content */}
      <div className="flex-1 p-8 bg-gray-50 overflow-auto text-gray-900">
        {activeTab === "Team" && !selectedConfig && (
          <div className="bg-white p-6 rounded-lg shadow">
            {/* UPDATED: Create button and text area for new config */}
            <div className="mb-4">
              <button
                className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
                onClick={() => setShowCreate(!showCreate)}
              >
                Create
              </button>
            </div>
            {showCreate && (
              <div className="mb-4">
                <textarea
                  className="w-full h-64 p-4 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                  placeholder="Enter new config JSON"
                  value={newConfig}
                  onChange={(e) => setNewConfig(e.target.value)}
                />
                {/* NEW: Input field for description */}
                <input
                  type="text"
                  placeholder="Enter configuration description"
                  value={newConfigDescription}
                  onChange={(e) => setNewConfigDescription(e.target.value)}
                  className="w-full p-3 mt-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <div className="flex justify-end mt-4">
                  <button
                    className="px-6 py-2 bg-green-600 text-white rounded-md shadow hover:bg-green-700 transition-colors"
                    onClick={() => {
                      fetch("http://localhost:8080/configs", {
                        method: "POST",
                        headers: {
                          "Content-Type": "application/json",
                          "Authorization": "Bearer " + jwtToken
                        },
                        body: JSON.stringify({ configContent: newConfig, description: newConfigDescription })
                      })
                        .then((res) => {
                          if (!res.ok) return res.text().then(text => { throw new Error(text) });
                          return res.json();
                        })
                        .then((data) => console.log("New config saved:", data))
                        .catch((err) => console.error("Error saving config:", err));
                      setNewConfig("");
                      setNewConfigDescription("");
                      setShowCreate(false);
                    }}
                  >
                    Save
                  </button>
                </div>
              </div>
            )}
            <h2 className="text-2xl font-bold mb-6">Team: {teamInfo?.name || "My Team"} Configurations</h2>
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-100">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Name</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Type</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Date</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Description</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {teamConfigs.map((config) => (
                  <tr key={config.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">{config.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.type}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.date}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <button
                        className="text-blue-600 hover:text-blue-800 font-medium transition-colors"
                        onClick={() => {
                          Promise.all([
                            fetch(`http://localhost:8080/configs/${config.id}`, {
                              headers: { "Authorization": "Bearer " + jwtToken }
                            }).then((res) => res.json()),
                            fetch(`http://localhost:8080/configs/${config.id}/changes`, {
                              headers: { "Authorization": "Bearer " + jwtToken }
                            }).then((res) => res.json()),
                          ])
                            .then(([cfgData, changes]) => {
                              setSelectedConfig(config.id);
                              setEditContent(JSON.stringify(cfgData, null, 2));
                              setConfigHistory(changes || []); // ensure changes is always an array
                              setIsEditing(false);
                            })
                            .catch((err) => console.error("Error fetching config details or changes:", err));
                        }}
                      >
                        Watch
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {activeTab === "Team" && selectedConfig && (
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold">
                {isEditing
                  ? `Editing Config ID: ${selectedConfig}`
                  : `Configuration ID: ${selectedConfig}`}
              </h2>
              <div className="space-x-3">
                {!isEditing && (
                  <button
                    className="text-blue-600 hover:text-blue-800 font-medium transition-colors"
                    onClick={() => setIsEditing(true)}
                  >
                    Edit
                  </button>
                )}
                <button
                  className="text-red-600 hover:text-red-800 font-medium transition-colors"
                  onClick={() => {
                    fetch(`http://localhost:8080/configs/${selectedConfig}/rollback`, {
                      method: "POST",
                      headers: { "Authorization": "Bearer " + jwtToken }
                    })
                      .then((res) => {
                        if (!res.ok) return res.text().then(text => { throw new Error(text) });
                        return res.json();
                      })
                      .then((data) => console.log("Rollback successful:", data))
                      .catch((err) => console.error("Error during rollback:", err));
                    setEditContent("");
                    setSelectedConfig(null);
                    setIsEditing(false);
                  }}
                >
                  Rollback
                </button>
              </div>
            </div>
            {isEditing && (
              <>
                <textarea
                  className="w-full h-64 p-4 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                />
                {/* Save button moved above change history */}
                <div className="flex justify-end mt-4">
                  <button
                    className="px-6 py-2 bg-blue-600 text-white rounded-md shadow hover:bg-blue-700 transition-colors"
                    onClick={() => {
                      fetch(`http://localhost:8080/configs/${selectedConfig}`, {
                        method: "PUT",
                        headers: {
                          "Content-Type": "application/json",
                          "Authorization": "Bearer " + jwtToken
                        },
                        body: JSON.stringify({ id: selectedConfig, configContent: editContent })
                      })
                        .then((res) => {
                          if (!res.ok)
                            return res.text().then(text => { throw new Error(text) });
                          return res.json();
                        })
                        .then((data) => console.log("Config updated:", data))
                        .catch((err) => console.error("Error updating config:", err));
                      setEditContent("");
                      setSelectedConfig(null);
                      setIsEditing(false);
                    }}
                  >
                    Save
                  </button>
                </div>
              </>
            )}
            {/* Always display change history */}
            <div className="mt-4">
              <h3 className="text-lg font-semibold mb-2">Change History</h3>
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                      Username
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                      Action
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                      Date
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {configHistory.map((change, idx) => (
                    <tr key={idx} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">{change.username}</td>
                      <td className="px-6 py-4 whitespace-nowrap">{change.action}</td>
                      <td className="px-6 py-4 whitespace-nowrap">{change.date}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {!isEditing && (
              <div className="text-xl text-gray-700 mt-4">
                Configuration content hidden. Press "Edit" to show content.
              </div>
            )}
          </div>
        )}

        {activeTab === "Admin" && (
          <div className="bg-gradient-to-br from-white to-gray-100 p-8 rounded-lg shadow-lg transition-all duration-300">
            <h2 className="text-3xl font-bold mb-6 text-center text-blue-600">Admin Panel</h2>
            {teamInfo ? (
              <>
                <h3 className="text-2xl font-semibold mb-4">Team: {teamInfo.name}</h3>
                {teamUsers.length > 0 ? (
                  <ul className="mb-4">
                    {teamUsers.map((user) => (
                      <li key={user.Id || user.Username} className="flex justify-between items-center p-2 border-b border-gray-200">
                        <span>{user.Username}</span>
                        <button
                          className="px-2 py-1 text-sm text-red-600 hover:text-red-800"
                          onClick={() => {
                            fetch("http://localhost:8080/teams/user", {
                              method: "PATCH",
                              headers: {
                                "Content-Type": "application/json",
                                "Authorization": "Bearer " + jwtToken
                              },
                              body: JSON.stringify({ username: user.Username })
                            })
                              .then((res) => {
                                if (!res.ok)
                                  return res.text().then(text => { throw new Error(text) });
                                // Handle empty response body gracefully:
                                return res.text().then(text => text ? JSON.parse(text) : {});
                              })
                              .then((data) => {
                                console.log("User removed:", data);
                                setTeamUsers(prev => prev.filter(u => u.Username !== user.Username));
                              })
                              .catch((err) => console.error("Error removing user:", err));
                          }}
                        >
                          Delete
                        </button>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <div className="text-gray-700 mb-4">No users found in your team.</div>
                )}
                {/* NEW: Add small window to add a user */}
                <div className="flex items-center space-x-3 mt-4">
                  <input
                    type="text"
                    placeholder="New Username"
                    value={newUserToAdd}
                    onChange={(e) => setNewUserToAdd(e.target.value)}
                    className="p-2 border border-gray-300 rounded"
                  />
                  <button
                    className="px-3 py-1 bg-green-600 text-white rounded"
                    onClick={() => {
                      fetch("http://localhost:8080/teams/user", {
                        method: "PATCH",
                        headers: {
                          "Content-Type": "application/json",
                          "Authorization": "Bearer " + jwtToken
                        },
                        body: JSON.stringify({ username: newUserToAdd, team_id: teamInfo.id })
                      })
                        .then((res) => res.text())
                        .then((text) => (text ? JSON.parse(text) : {}))
                        .then((data) => {
                          console.log("User added:", data);
                          setTeamUsers(prev => [...prev, { Username: newUserToAdd }]);
                          setNewUserToAdd("");
                        })
                        .catch((err) => console.error("Error adding user:", err));
                    }}
                  >
                    Add User
                  </button>
                </div>
              </>
            ) : (
              <div className="flex flex-col items-center">
                <div className="text-gray-700 mb-4">No team found. Please create a team.</div>
                <div className="flex items-center space-x-3">
                  <input
                    type="text"
                    placeholder="Team Name"
                    value={newTeamName}
                    onChange={(e) => setNewTeamName(e.target.value)}
                    className="flex-1 p-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <button
                    className="px-5 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                    onClick={() => {
                      fetch("http://localhost:8080/teams", {
                        method: "POST",
                        headers: {
                          "Content-Type": "application/json",
                          "Authorization": "Bearer " + jwtToken
                        },
                        body: JSON.stringify({ team: { name: newTeamName } })
                      })
                        .then((res) => {
                          if (!res.ok) return res.text().then(text => { throw new Error(text) });
                          return res.json();
                        })
                        .then((data) => {
                          console.log("Team created:", data);
                          setTeamInfo(data);
                        })
                        .catch((err) => console.error(err));
                      setNewTeamName("");
                    }}
                  >
                    Create Team
                  </button>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === "Info" && (
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-4">Info</h2>
            <p className="text-gray-900">Some information or tips can be displayed here.</p>
          </div>
        )}

        {activeTab === "Login/Register" && (
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-4">{authMode === "login" ? "Login" : "Register"}</h2>
            {/* Switch button */}
            <div className="mb-4 flex justify-end">
              <button
                className="text-sm text-blue-600 hover:underline"
                onClick={() => setAuthMode(authMode === "login" ? "register" : "login")}
              >
                Switch to {authMode === "login" ? "Register" : "Login"}
              </button>
            </div>
            {authMode === "login" ? (
              <div className="mb-6">
                <input
                  type="text"
                  placeholder="Username"
                  value={loginUsername}
                  onChange={(e) => setLoginUsername(e.target.value)}
                  className="w-full p-3 mb-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <input
                  type="password"
                  placeholder="Password"
                  value={loginPassword}
                  onChange={(e) => setLoginPassword(e.target.value)}
                  className="w-full p-3 mb-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                  onClick={handleLogin}
                >
                  Login
                </button>
              </div>
            ) : (
              <div>
                <input
                  type="text"
                  placeholder="Username"
                  value={registerUsername}
                  onChange={(e) => setRegisterUsername(e.target.value)}
                  className="w-full p-3 mb-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <input
                  type="email"
                  placeholder="Email"
                  value={registerEmail}
                  onChange={(e) => setRegisterEmail(e.target.value)}
                  className="w-full p-3 mb-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <input
                  type="password"
                  placeholder="Password"
                  value={registerPassword}
                  onChange={(e) => setRegisterPassword(e.target.value)}
                  className="w-full p-3 mb-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                  className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
                  onClick={handleRegister}
                >
                  Register
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}