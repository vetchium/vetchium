import { useState, useEffect, useRef } from "react";
import { config } from "@/config";
import Cookies from "js-cookie";
import { AddOrgUserRequest, OrgUser } from "@psankar/vetchi-typespec";

export function useOrgUsers() {
  const [users, setUsers] = useState<OrgUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchUsersRef = useRef(async (includeDisabled: boolean = false) => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/filter-org-users`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({
            state: includeDisabled
              ? ["ACTIVE_ORG_USER", "ADDED_ORG_USER", "DISABLED_ORG_USER"]
              : ["ACTIVE_ORG_USER", "ADDED_ORG_USER"],
          }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to fetch organization users");
      }

      const data = await response.json();
      setUsers(data);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error
          ? err
          : new Error("Failed to fetch organization users")
      );
    } finally {
      setIsLoading(false);
    }
  });

  useEffect(() => {
    fetchUsersRef.current();
  }, []);

  const addUser = async (data: AddOrgUserRequest) => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-org-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to add organization user");
      }

      await fetchUsersRef.current(); // Use the ref here
    } catch (err) {
      throw err instanceof Error
        ? err
        : new Error("Failed to add organization user");
    }
  };

  const disableUser = async (email: string) => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/disable-org-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({ email }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to disable organization user");
      }

      await fetchUsersRef.current(); // Use the ref here
    } catch (err) {
      throw err instanceof Error
        ? err
        : new Error("Failed to disable organization user");
    }
  };

  const enableUser = async (email: string) => {
    try {
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/enable-org-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${Cookies.get("session_token")}`,
          },
          body: JSON.stringify({ email }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to enable organization user");
      }

      await fetchUsersRef.current(); // Use the ref here
    } catch (err) {
      throw err instanceof Error
        ? err
        : new Error("Failed to enable organization user");
    }
  };

  return {
    users,
    isLoading,
    error,
    addUser,
    disableUser,
    enableUser,
    fetchUsers: fetchUsersRef.current,
  };
}
