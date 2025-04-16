import { config } from "@/config";
import {
  FollowStatus,
  FollowUserRequest,
  GetFollowStatusRequest,
  UnfollowUserRequest,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslation } from "./useTranslation";

interface UseFollowUserResult {
  followStatus: FollowStatus | null;
  isLoadingStatus: boolean;
  isFollowing: boolean;
  isUnfollowing: boolean;
  error: Error | null;
  getFollowStatus: (handle: string) => Promise<void>;
  followUser: (handle: string) => Promise<void>;
  unfollowUser: (handle: string) => Promise<void>;
  clearError: () => void;
}

export function useFollowUser(): UseFollowUserResult {
  const router = useRouter();
  const { t } = useTranslation();
  const [followStatus, setFollowStatus] = useState<FollowStatus | null>(null);
  const [isLoadingStatus, setIsLoadingStatus] = useState<boolean>(false);
  const [isFollowing, setIsFollowing] = useState<boolean>(false);
  const [isUnfollowing, setIsUnfollowing] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  const getFollowStatus = async (handle: string) => {
    try {
      setIsLoadingStatus(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: GetFollowStatusRequest = { handle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-follow-status`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(t("profile.error.userNotFound"));
        }
        throw new Error(t("profile.error.followStatusFailed"));
      }

      const data = await response.json();
      setFollowStatus(data);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsLoadingStatus(false);
    }
  };

  const followUser = async (handle: string) => {
    try {
      setIsFollowing(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: FollowUserRequest = { handle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/follow-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(t("profile.error.userNotFound"));
        }
        throw new Error(t("profile.error.followFailed"));
      }

      // After successful follow, update the status
      await getFollowStatus(handle);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsFollowing(false);
    }
  };

  const unfollowUser = async (handle: string) => {
    try {
      setIsUnfollowing(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: UnfollowUserRequest = { handle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/unfollow-user`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(t("profile.error.userNotFound"));
        }
        throw new Error(t("profile.error.unfollowFailed"));
      }

      // After successful unfollow, update the status
      await getFollowStatus(handle);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error(t("common.error.serverError"))
      );
      throw err;
    } finally {
      setIsUnfollowing(false);
    }
  };

  const clearError = () => {
    setError(null);
  };

  return {
    followStatus,
    isLoadingStatus,
    isFollowing,
    isUnfollowing,
    error,
    getFollowStatus,
    followUser,
    unfollowUser,
    clearError,
  };
}
