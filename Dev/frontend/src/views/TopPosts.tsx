import { useEffect, useState } from "react";
import Post from "../components/Post";
import { Post as PostType } from "../utils/types"
import { getPosts, getTopPostIds } from "../utils/serverFunctions";

function TopPosts() {
  const [postLimit, setPostLimit] = useState(1)
  const [posts, setPosts] = useState<PostType[]>([])
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    setLoading(true)
    getTopPostIds(postLimit).then(idsRes => {
      getPosts(idsRes.Result).then((postsRes) => {
        setLoading(false)
        if (postsRes.Success) {
          setPosts(postsRes.Result)
        }
      })
    })
  }, [postLimit])

  return <>
    {posts.map((p) => <Post loading={false} src={p.Url} caption={p.Caption} random={false} comments={true} />)}
  </>

}

export default TopPosts
