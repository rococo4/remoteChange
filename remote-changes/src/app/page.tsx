"use client";

import { useState } from "react";

export default function Page() {
  const [activeTab, setActiveTab] = useState("Team");
  const [selectedConfig, setSelectedConfig] = useState<keyof typeof configDetails | null>(null);
  const [editContent, setEditContent] = useState("");

  // New states for admin team management
  const [adminTeam, setAdminTeam] = useState<{ name: string; users: string[] } | null>(null);
  const [newTeamName, setNewTeamName] = useState("");
  const [newUserName, setNewUserName] = useState("");

  const teamConfigs: { name: keyof typeof configDetails; type: string; date: string; description: string }[] = [
    { name: "Config 1", type: "Type A", date: "2023-10-01", description: "Description 1" },
    { name: "Config 2", type: "Type B", date: "2023-10-02", description: "Description 2" },
  ];

  const configDetails = {
    "Config 1": [
      { username: "User1", type: "Type A", date: "2023-10-01", description: "User Description 1" },
    ],
    "Config 2": [
      { username: "User2", type: "Type B", date: "2023-10-02", description: "User Description 2" },
    ],
  };

  return (
    <div className="flex h-screen font-sans">
      {/* Sidebar */}
      <div className="w-64 bg-gradient-to-b from-blue-600 to-blue-800 text-white p-6 shadow-lg">
        <ul className="space-y-6">
          {["Team", "Admin", "Info"].map((tab) => (
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
          ))}
        </ul>
      </div>

      {/* Main Content */}
      <div className="flex-1 p-8 bg-gray-50 overflow-auto text-gray-900">
        {activeTab === "Team" && !selectedConfig && (
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-6">Team Configurations</h2>
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
                  <tr key={config.name} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">{config.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.type}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.date}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{config.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <button
                        className="text-blue-600 hover:text-blue-800 font-medium transition-colors"
                        onClick={() => setSelectedConfig(config.name)}
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
              <h2 className="text-2xl font-bold">{selectedConfig} Configuration</h2>
              <button
                className="text-blue-600 hover:text-blue-800 font-medium transition-colors"
                onClick={() =>
                  setEditContent(JSON.stringify(configDetails[selectedConfig], null, 2))
                }
              >
                Edit
              </button>
            </div>
            {!editContent ? (
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Username</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Type</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Date</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">Description</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {configDetails[selectedConfig].map((detail, index) => (
                    <tr key={index} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">{detail.username}</td>
                      <td className="px-6 py-4 whitespace-nowrap">{detail.type}</td>
                      <td className="px-6 py-4 whitespace-nowrap">{detail.date}</td>
                      <td className="px-6 py-4 whitespace-nowrap">{detail.description}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div>
                <textarea
                  className="w-full h-64 p-4 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                />
                <div className="flex justify-end mt-4">
                  <button
                    className="px-6 py-2 bg-blue-600 text-white rounded-md shadow hover:bg-blue-700 transition-colors"
                    onClick={() => {
                      // Implement save functionality here
                      console.log("Saved content:", editContent);
                      setEditContent("");
                    }}
                  >
                    Save
                  </button>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === "Admin" && (
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-4">Admin Panel</h2>
            {adminTeam ? (
              <>
                {/* Display team users */}
                <h3 className="text-xl font-semibold mb-2">Team: {adminTeam.name}</h3>
                <ul className="mb-4">
                  {adminTeam.users.map((username, index) => (
                    <li key={index} className="flex justify-between items-center py-1">
                      <span>{username}</span>
                      <button
                        className="text-red-600 hover:text-red-800 font-medium"
                        onClick={() =>
                          setAdminTeam({
                            ...adminTeam,
                            users: adminTeam.users.filter((_, i) => i !== index),
                          })
                        }
                      >
                        Delete
                      </button>
                    </li>
                  ))}
                </ul>
                {/* Add new user form */}
                <div className="flex items-center space-x-2">
                  <input
                    type="text"
                    placeholder="Username"
                    value={newUserName}
                    onChange={(e) => setNewUserName(e.target.value)}
                    className="p-2 border rounded"
                  />
                  <button
                    className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                    onClick={() => {
                      if (newUserName.trim() === "") return;
                      setAdminTeam({ ...adminTeam, users: [...adminTeam.users, newUserName.trim()] });
                      setNewUserName("");
                    }}
                  >
                    Add User
                  </button>
                </div>
              </>
            ) : (
              // If no team exists then display create team form
              <div className="flex items-center space-x-2">
                <input
                  type="text"
                  placeholder="Team Name"
                  value={newTeamName}
                  onChange={(e) => setNewTeamName(e.target.value)}
                  className="p-2 border rounded"
                />
                <button
                  className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                  onClick={() => {
                    if (newTeamName.trim() === "") return;
                    setAdminTeam({ name: newTeamName.trim(), users: [] });
                    setNewTeamName("");
                  }}
                >
                  Create Team
                </button>
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
      </div>
    </div>
  );
}