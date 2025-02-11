import { useState, useEffect } from "react";
import {
  WorkHistory as WorkHistoryType,
  AddWorkHistoryRequest,
  UpdateWorkHistoryRequest,
  DeleteWorkHistoryRequest,
  ListWorkHistoryRequest,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useTranslation } from "@/hooks/useTranslation";

interface WorkHistoryProps {
  userHandle: string;
  canEdit: boolean;
}

type WorkHistoryFormData = Omit<AddWorkHistoryRequest, "id">;

export function WorkHistory({ userHandle, canEdit }: WorkHistoryProps) {
  const { t } = useTranslation();
  const [workHistory, setWorkHistory] = useState<WorkHistoryType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState<string | null>(null);
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [formData, setFormData] = useState<WorkHistoryFormData>({
    employer_domain: "",
    title: "",
    start_date: "",
  });

  useEffect(() => {
    fetchWorkHistory();
  }, [userHandle]);

  async function fetchWorkHistory() {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        // Handle unauthenticated state
        return;
      }

      const request: ListWorkHistoryRequest = { user_handle: userHandle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/list-work-history`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) throw new Error(t("workHistory.error.fetchFailed"));

      const data = await response.json();
      setWorkHistory(data);
    } catch (error) {
      console.error("Error fetching work history:", error);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        // Handle unauthenticated state
        return;
      }

      const endpoint = isEditing
        ? `${config.API_SERVER_PREFIX}/hub/update-work-history`
        : `${config.API_SERVER_PREFIX}/hub/add-work-history`;

      const method = isEditing ? "PUT" : "POST";
      const body = isEditing
        ? ({ ...formData, id: isEditing } as UpdateWorkHistoryRequest)
        : (formData as AddWorkHistoryRequest);

      const response = await fetch(endpoint, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) throw new Error(t("workHistory.error.saveFailed"));

      await fetchWorkHistory();

      setFormData({
        employer_domain: "",
        title: "",
        start_date: "",
      });
      setIsEditing(null);
      setIsAddingNew(false);
    } catch (error) {
      console.error("Error saving work history:", error);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm(t("workHistory.deleteConfirm"))) {
      return;
    }

    try {
      const token = Cookies.get("session_token");
      if (!token) {
        // Handle unauthenticated state
        return;
      }

      const request: DeleteWorkHistoryRequest = { id };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-work-history`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) throw new Error(t("workHistory.error.deleteFailed"));

      await fetchWorkHistory();
    } catch (error) {
      console.error("Error deleting work history:", error);
    }
  }

  if (isLoading) {
    return <div>{t("workHistory.loading")}</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-semibold">{t("workHistory.title")}</h2>
        {canEdit && !isAddingNew && !isEditing && (
          <button
            onClick={() => setIsAddingNew(true)}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            {t("workHistory.addExperience")}
          </button>
        )}
      </div>

      {(isAddingNew || isEditing) && canEdit && (
        <form onSubmit={handleSubmit} className="mb-8 space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              {t("workHistory.companyDomain")}
            </label>
            <input
              type="text"
              value={formData.employer_domain}
              onChange={(e) =>
                setFormData({ ...formData, employer_domain: e.target.value })
              }
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              {t("workHistory.jobTitle")}
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) =>
                setFormData({ ...formData, title: e.target.value })
              }
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              {t("workHistory.startDate")}
            </label>
            <input
              type="date"
              value={formData.start_date}
              onChange={(e) =>
                setFormData({ ...formData, start_date: e.target.value })
              }
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              {t("workHistory.endDate")}
            </label>
            <input
              type="date"
              value={formData.end_date || ""}
              onChange={(e) =>
                setFormData({ ...formData, end_date: e.target.value })
              }
              className="w-full p-2 border rounded"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              {t("workHistory.description")}
            </label>
            <textarea
              value={formData.description || ""}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              className="w-full p-2 border rounded"
              rows={4}
            />
          </div>
          <div className="flex space-x-4">
            <button
              type="submit"
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            >
              {isEditing
                ? t("workHistory.actions.save")
                : t("workHistory.addExperience")}
            </button>
            <button
              type="button"
              onClick={() => {
                setIsEditing(null);
                setIsAddingNew(false);
                setFormData({
                  employer_domain: "",
                  title: "",
                  start_date: "",
                });
              }}
              className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600"
            >
              {t("workHistory.actions.cancel")}
            </button>
          </div>
        </form>
      )}

      <div className="space-y-6">
        {workHistory.map((entry) => (
          <div key={entry.id} className="border rounded p-4">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="font-semibold">{entry.title}</h3>
                <p className="text-gray-600">
                  {entry.employer_name || entry.employer_domain}
                </p>
                <p className="text-sm text-gray-500">
                  {new Date(entry.start_date).toLocaleDateString()} -{" "}
                  {entry.end_date
                    ? new Date(entry.end_date).toLocaleDateString()
                    : t("workHistory.present")}
                </p>
                {entry.description && (
                  <p className="mt-2 text-gray-700">{entry.description}</p>
                )}
              </div>
              {canEdit && !isEditing && !isAddingNew && (
                <div className="space-x-2">
                  <button
                    onClick={() => {
                      setIsEditing(entry.id);
                      setFormData({
                        employer_domain: entry.employer_domain,
                        title: entry.title,
                        start_date: entry.start_date,
                        end_date: entry.end_date,
                        description: entry.description,
                      });
                    }}
                    className="text-blue-500 hover:text-blue-600"
                  >
                    {t("workHistory.actions.edit")}
                  </button>
                  <button
                    onClick={() => handleDelete(entry.id)}
                    className="text-red-500 hover:text-red-600"
                  >
                    {t("workHistory.actions.delete")}
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
