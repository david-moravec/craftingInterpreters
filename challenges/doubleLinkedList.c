#include<stddef.h>
#include<stdlib.h>
#include<stdio.h>
#include<assert.h>
#include<stdbool.h>

//Should be generic as well
typedef int Content;

//somethign that will be generic in the future, dont know how to declare generic type
typedef float T;

typedef struct _Node {
    Content content;
    struct _Node* next;
    struct _Node* prev;
} Node;

typedef struct _LinkedList {
    Node* head;
    Node* tail;
} LinkedList;


void append_to(Node* parent, Content val) {
    Node* child = (Node*) malloc(sizeof(Node));
    child->content = val;
    child->prev = parent;

    if (parent->next == NULL){
        parent->next = child;
    //insert child behind parent append what was previously appended to parent
    } else { 
        Node* parent_next = parent->next;
        parent->next = child;
        child->next = parent_next;
        parent_next->prev = child;
    }
}

bool append_nodes(Node* parent, Node* to_append) {
    //make sure parent does not have next first
    if (parent->next != NULL) {
        return false;
    } else {
        parent->next = to_append;
        return true;
    }
}

bool delete_node(Node* current) {
    Node* parent = current->prev;
    Node* child = current->next;

    parent->next = NULL;

    if (append_nodes(parent, child)) {
        free(current);
        return true;
    } else {
        return false;
    } 
}

Node* head_linked_list(Node* current) {
    if (current->prev == NULL) {
        return current;
    }

    return head_linked_list(current->prev);
}

Node* tail_linked_list(Node* current) {
    if (current->next == NULL) {
        return current;
    }

    return tail_linked_list(current->next);
}

typedef T (*FUNC)(Content, T)


T foldr(T (*func)(Content, T), T acc, LinkedList list) {
    Node* next = list.head;

    if (next != list.tail) {
        func(next->content, acc);
    }
}



void delete(Node* current, Content to_delete) {
    //Deletes all nodes with content == to_delete
    //Delete forward
    //
    //

    Node* start = current;

    Node* head = head_linked_list(current);
    Node* tail = tail_linked_list(current);


    while (current != tail) {
        Node* next = current->next;

        if (current->content == to_delete) {
            Node* prev = current->prev;
            prev->next = next;
            free(current);
        }

        current = next;
    }



    while (current != head) {
        Node* prev = current->prev;

        if (current->content == to_delete) {
            Node* next = current->next;
            next->prev = prev;
            free(current);
        }

        current = prev;
    }
}




int len_linked_list(Node* some_node) {
        //start from head
        Node* current = head_linked_list(some_node);
        Node* tail = tail_linked_list(some_node);

        int sum = 1;

        while (current != tail) {
            sum += 1;

            current = current->next;
        }

        return sum;
    }


int main() {
    Node* head = NULL;
    head = (Node*) malloc(sizeof(Node));

    if (head == NULL) {
        printf("Head still NULL");
        return 1;
    }

    head->content = 2;
    assert(head->content == 2);

    append_to(head, 4);
    assert(head->next->content == 4);

    append_to(head, 3);
    assert(head->next->content == 3);
    assert(head->next->next->content == 4);

    assert(head == head_linked_list(head->next->next));
    assert(head->next->next == tail_linked_list(head));

    assert(len_linked_list(head) == 3);

        
    if (delete_node(head->next)) {
        assert(head->next->content == 4);
        assert(len_linked_list(head) == 2);
    };

    for (int i = 0; i<10; i ++) {
        append_to(head, i);
    }

    assert(len_linked_list(head) == 12);

    return 0;
}
