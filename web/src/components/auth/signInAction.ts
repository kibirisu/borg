import type { ActionFunctionArgs } from "react-router";
import { redirect } from "react-router";
import type { AppClient } from "../../lib/client";

export function signInAction(client: AppClient) {
  return async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();

    const username = formData.get("username")?.toString();
    const password = formData.get("password")?.toString();

    const errors: Record<string, string> = {};

    if (!username) {
      errors.username = "Field is mandatory";
    } else if (username.length < 6) {
      errors.username = "Username should be at least 6 characters";
    }

    if (!password) {
      errors.password = "Field is mandatory";
    } else if (password.length < 6) {
      errors.password = "Password should be at least 6 characters";
    }

    if (Object.keys(errors).length > 0) {
      return errors;
    }

    const safeUsername = username!;
    const safePassword = password!;

    const mutation = async () => {
      return client.fetchClient.POST("/api/auth/login", {
        body: { username: safeUsername, password: safePassword },
      });
    };

    try {
      await mutation();
    } catch {
      return redirect("/?alert=user-missing");
    }

    return redirect("/explore");
  };
}
