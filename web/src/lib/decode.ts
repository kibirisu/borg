export interface DecodedToken {
  username: string | null;
  userId: number | null;
}

const decodeToken = (token: string | null): DecodedToken | null => {
  if (!token) {
    return null;
  }
  const raw = token.replace(/^Bearer:\s*/i, "");
  const base64Url = raw.split(".")[1];
  if (!base64Url) {
    return null;
  }
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join(""),
  );
  const payload = JSON.parse(jsonPayload) as {
    username?: string;
    userId?: number | string;
  };
  return {
    username: payload.username ?? null,
    userId:
      typeof payload.userId === "number"
        ? payload.userId
        : payload.userId
          ? Number(payload.userId)
          : null,
  };
};

export default decodeToken;
