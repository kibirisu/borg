import { type ActionFunctionArgs, redirect } from "react-router";
import type { AppClient } from "../../lib/client";

const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/;

export function signUpAction(client: AppClient) {
  return async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();

    const username = formData.get("username")?.toString();
    const password = formData.get("password")?.toString();
    const confirmPassword = formData.get("confirmPassword")?.toString();

    const errors: Record<string, string> = {};

    if (!username) {
      errors.username = "Field is mandatory";
    } else if (!USERNAME_REGEX.test(username)) {
      errors.username =
        "Username can only contain letters, numbers and underscores.";
    } else if (username.length > 30) {
      errors.username = "Username is too long (max 30 chars).";
    }

    if (!password) {
      errors.password = "Field is mandatory";
    } else if (password.length < 1) {
      errors.password = "Password should be at least 6 characters";
    }

    if (!confirmPassword) {
      errors.confirmPassword = "Field is mandatory";
    }

    if (password && confirmPassword && password !== confirmPassword) {
      errors.confirmPassword = "Passwords do not match";
    }

    if (Object.keys(errors).length > 0) {
      return errors;
    }

    if (!username || !password) {
      return errors;
    }

    const safeUsername = username;
    const safePassword = password;

    const mutation = async () => {
      return client.fetchClient.POST("/auth/register", {
        body: { username: safeUsername, password: safePassword },
      });
    };

    try {
      const res = await mutation();
      if (res.error) {
        console.error("[signup] api error");
        return { form: "Registration failed" };
      }
    } catch (err) {
      console.error("[signup] network/client error", err);
      return { form: "Registration failed" };
    }

    return redirect("/signin");
  };
}
